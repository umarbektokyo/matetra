package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"matetra/api"
	"matetra/model"
	"matetra/utils"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	log.SetFlags(0)

	cmd := os.Args

	if len(cmd) < 2 {
		fmt.Println("usage: go run matetra <server-address>")
		fmt.Println("example: go run matetra localhost:1729")
		return
	}

	serverAddr := cmd[1]

	if !strings.HasPrefix(serverAddr, "ws://") {
		serverAddr = "ws://" + serverAddr
	}

	u, err := url.Parse(serverAddr)
	if err != nil {
		log.Fatalf("error: invalid server address format: %v", err)
	}
	u.Path = "/ws"

	fmt.Printf("attempting to connect to server at %s...\n", serverAddr)

	// connect to the websocket server
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("error: could not connect to server at %s. Is the server still running? %v", u.String(), err)
	}
	defer c.Close()
	fmt.Println("Connection successful.")

	// register the player
	if err := registerPlayer(c); err != nil {
		log.Fatalf("registration failed: %v", err)
	}

	// start listening for server updates
	listenForUpdates(c)
}

func registerPlayer(c *websocket.Conn) error {
	reader := bufio.NewReader(os.Stdin)

	// get username
	fmt.Print("username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// get password
	fmt.Print("password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	passwordHash := utils.Hash(password)
	fmt.Printf("hashed password (SHA256): %s...\n", passwordHash[:8])

	// construct the PlayerPayload and the ADD_PLAYER message
	payload := api.PlayerPayload{
		Name: username,
		Hash: passwordHash,
	}
	addPlayerMsg := api.Message{
		Type:    "ADD_PLAYER",
		Payload: payload,
	}

	// send ADD_PLAYER request
	fmt.Println("registering player...")
	if err := c.WriteJSON(addPlayerMsg); err != nil {
		return fmt.Errorf("error sending registration request: %v", err)
	}

	// wait for server's response
	var response api.Message
	if err := c.ReadJSON(&response); err != nil {
		return fmt.Errorf("error reading registration response: %v", err)
	}

	switch response.Type {
	case "PLAYER_ADDED":
		fmt.Printf("success: player @%s has been added to the game!\n", username)
		return nil
	case "ERROR":
		// handle the case where the server explicitly sends an error message
		errorPayloadBytes, err := json.Marshal(response.Payload)
		if err != nil {
			return fmt.Errorf("registration failed: unknown error format")
		}
		var errorData map[string]string
		json.Unmarshal(errorPayloadBytes, &errorData)
		return fmt.Errorf("registration failed: %s", errorData["message"])
	default:
		return fmt.Errorf("unexpected server response type: %s", response.Type)
	}
}

func listenForUpdates(c *websocket.Conn) {
	fmt.Println("Listening for game updates...")
	for {
		var msg api.Message
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Println("Server connection closed.")
				return
			}
			log.Printf("Error reading message: %v", err)
			return
		}

		switch msg.Type {
		case "STATE_UPDATE":
			// the state is in the payload, we need to unmarshal it into the GameState model
			statePayloadBytes, err := json.Marshal(msg.Payload)
			if err != nil {
				log.Printf("Error marshalling state payload: %v", err)
				continue
			}
			var gameState model.GameState
			if err := json.Unmarshal(statePayloadBytes, &gameState); err != nil {
				log.Printf("Error unmarshalling GameState: %v", err)
				continue
			}
			fmt.Printf("\n--- Game State Update (Turn: %d) ---\n", gameState.Turn)
			fmt.Printf("Players: %d, Cards in Deck: %d\n", len(gameState.Players), countDeckCards(gameState))
			fmt.Println("---------------------------------")
		case "ERROR":
			// handle server-pushed errors
			errorPayloadBytes, _ := json.Marshal(msg.Payload)
			var errorData map[string]string
			json.Unmarshal(errorPayloadBytes, &errorData)
			fmt.Printf("\n[SERVER ERROR]: %s\n", errorData["message"])
		default:
			fmt.Printf("Received unhandled message type: %s\n", msg.Type)
		}
	}
}

func countDeckCards(gs model.GameState) int {
	count := 0
	for _, card := range gs.Cards {
		if card.Owner == -1 {
			count++
		}
	}
	return count
}
