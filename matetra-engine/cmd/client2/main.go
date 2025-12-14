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
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// Global state tracking variables
var CurrentGameState model.GameState
var PlayerID int = -1
var PlayerName string

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

	fmt.Printf("Attempting to connect to server at %s...\n", serverAddr)

	// connect to the websocket server
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("error: could not connect to server at %s. Is the server still running? %v", u.String(), err)
	}
	defer c.Close()
	fmt.Println("Connection successful.")

	// register the player
	if err := registerPlayer(c); err != nil {
		log.Fatalf("Registration failed: %v", err)
	}

	// Use goroutines for simultaneous listening and commanding
	go listenForUpdates(c)
	commandLoop(c)
}

func registerPlayer(c *websocket.Conn) error {
	reader := bufio.NewReader(os.Stdin)

	// get username
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// get password
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	passwordHash := utils.Hash(password)
	fmt.Printf("Hashed password (SHA256): %s...\n", passwordHash[:8])

	payload := api.PlayerPayload{
		Name: username,
		Hash: passwordHash,
	}
	addPlayerMsg := api.Message{
		Type:    "ADD_PLAYER",
		Payload: payload,
	}

	fmt.Println("Registering player...")
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
		// Find the PlayerID from the received state update
		// NOTE: The server should really return the PlayerID in the PLAYER_ADDED payload for simplicity.
		// For now, we rely on the first STATE_UPDATE to follow.

		PlayerName = username
		// We can't safely set PlayerID here, rely on STATE_UPDATE after the loop starts.
		fmt.Printf("Success: player @%s has been added to the game!\n", PlayerName)
		return nil
	case "ERROR":
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

// ----------------------------------------------------------------------
// GAME STATE DISPLAY
// ----------------------------------------------------------------------

func displayGameState(gs model.GameState) {
	fmt.Print("\033[H\033[2J") // Clear terminal screen

	// Set PlayerID globally if not yet set
	if PlayerID == -1 {
		for i, p := range gs.Players {
			if p.Name == PlayerName {
				PlayerID = i
				break
			}
		}
	}

	fmt.Println("=====================================================================")
	fmt.Printf(" 🎮 GAME: %s | TURN: %d | CURRENT PLAYER: @%s\n", gs.GameID, gs.Turn, gs.Players[gs.Turn%len(gs.Players)].Name)
	fmt.Println("=====================================================================")

	// 1. Display Player Numbers
	fmt.Println("\n--- PLAYER NUMBERS ---")
	for i, p := range gs.Players {
		doneStatus := "✅ DONE"
		if i < len(gs.Done) && !gs.Done[i] {
			doneStatus = "▶️ ACTIVE"
		}

		marker := "  "
		if i == PlayerID {
			marker = ">>"
		} else if i == (gs.Turn % len(gs.Players)) {
			marker = "🎯"
		}

		numberStrings := make([]string, 5)
		for j, num := range gs.Numbers[i] {
			numberStrings[j] = fmt.Sprintf("[%d:%s%s]", j, num.Value, num.Mark)
		}

		fmt.Printf("%s %s (@%s) [%s]: %s\n", marker, doneStatus, p.Name, strings.Join(numberStrings, " | "), p.Hash[:8])
	}

	// 2. Display Player Hand
	fmt.Println("\n--- YOUR HAND ---")
	handCount := 0
	for i, card := range gs.Cards {
		if card.Owner == PlayerID {
			handCount++
			// Find the required input string from the card
			inputsReq := ""
			if i < len(gs.Cards) {
				inputsReq = gs.Cards[i].InputsReq
			}
			fmt.Printf("  [C:%d] %s (Req: %s) -> %s\n", i, card.Name, inputsReq, card.Description)
		}
	}
	if handCount == 0 {
		fmt.Println("  (Your hand is empty)")
	}

	// 3. Display Queue
	fmt.Println("\n--- MOVE QUEUE ---")
	if len(gs.Queue) > 0 {
		queueDetails := make([]string, len(gs.Queue))
		for i, cardIndex := range gs.Queue {
			// Find the card name
			cardName := "Unknown Card"
			if cardIndex >= 0 && cardIndex < len(gs.Cards) {
				cardName = gs.Cards[cardIndex].Name
			}
			queueDetails[i] = fmt.Sprintf("%s (ID:%d)", cardName, cardIndex)
		}
		fmt.Printf("  %s\n", strings.Join(queueDetails, " -> "))
	} else {
		fmt.Println("  (Queue is empty)")
	}

	fmt.Println("---------------------------------------------------------------------")
}

// ----------------------------------------------------------------------
// MESSAGE LISTENER
// ----------------------------------------------------------------------

func listenForUpdates(c *websocket.Conn) {
	for {
		var msg api.Message
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Println("\nServer connection closed.")
				os.Exit(0)
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		switch msg.Type {
		case "PLAY_CARD_REPLY":
			// Both the success/fail confirmation and the permanent broadcast use this.
			var reply api.CardPlayReply
			payloadBytes, _ := json.Marshal(msg.Payload)
			if err := json.Unmarshal(payloadBytes, &reply); err != nil {
				log.Printf("Error unmarshalling CardPlayReply: %v", err)
				continue
			}

			// 1. Update Global State
			if reply.NewGameState != nil {
				CurrentGameState = *reply.NewGameState
				displayGameState(CurrentGameState)
			}

			// 2. Display Message
			prefix := "[INFO]"
			if !reply.Success {
				prefix = "[ERROR]"
			}
			fmt.Printf("\n%s %s\n", prefix, reply.Message)

		case "ERROR":
			errorPayloadBytes, _ := json.Marshal(msg.Payload)
			var errorData map[string]string
			json.Unmarshal(errorPayloadBytes, &errorData)
			fmt.Printf("\n[SERVER ERROR]: %s\n", errorData["message"])
		default:
			// Ignore unhandled types like "PLAYER_ADDED" since we handle it in registerPlayer
			// and rely on subsequent PLAY_CARD_REPLY or STATE_UPDATE for state changes.
		}
	}
}

// ----------------------------------------------------------------------
// COMMAND INTERFACE
// ----------------------------------------------------------------------

func commandLoop(c *websocket.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Ensure the command prompt appears clearly after the state
		fmt.Printf("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.FieldsFunc(input, func(r rune) bool {
			return r == '(' || r == ')' || r == ',' || r == ' '
		})

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "apply":
			if len(parts) < 3 {
				fmt.Println("Usage: apply(cardIndex, input1, input2, ..., permanent)")
				fmt.Println("Permanent: 1 for yes, 0 for preview. Example: apply(2, 0, 1, 1)")
				continue
			}

			// Parse card index
			cardIndex, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid card index.")
				continue
			}

			// Extract inputs and 'permanent' flag
			var inputs []int
			for i := 2; i < len(parts)-1; i++ {
				inputVal, err := strconv.Atoi(parts[i])
				if err != nil {
					fmt.Printf("Invalid input value at position %d.\n", i-1)
					inputs = nil
					break
				}
				inputs = append(inputs, inputVal)
			}
			if inputs == nil && len(parts) > 2 {
				continue // Error already reported
			}

			// Parse permanent flag
			permanentInt, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil || (permanentInt != 0 && permanentInt != 1) {
				fmt.Println("Invalid permanent flag. Use 1 for permanent, 0 for preview.")
				continue
			}

			permanent := permanentInt == 1

			sendPlayCard(c, cardIndex, inputs, permanent)

		case "turnend":
			sendTurnEnd(c)

		case "exit", "quit":
			fmt.Println("Exiting client.")
			return

		case "state":
			if CurrentGameState.GameID != "" {
				displayGameState(CurrentGameState)
			} else {
				fmt.Println("Waiting for initial game state...")
			}
		default:
			fmt.Printf("Unknown command: %s. Use apply(), turnend(), or exit().\n", command)
		}
	}
}

// ----------------------------------------------------------------------
// COMMAND SENDERS
// ----------------------------------------------------------------------

func sendPlayCard(c *websocket.Conn, cardIndex int, inputs []int, permanent bool) {
	if PlayerID == -1 {
		fmt.Println("[ERROR] Player ID not yet established. Cannot move.")
		return
	}

	playCardMsg := api.Message{
		Type: "PLAY_CARD",
		Payload: api.CardPlayPayload{
			CardIndex: cardIndex,
			Inputs:    inputs,
			Permanent: permanent,
		},
	}

	if err := c.WriteJSON(playCardMsg); err != nil {
		log.Printf("Error sending PLAY_CARD: %v", err)
	}
	if !permanent {
		fmt.Println("Sent preview request...")
	} else {
		fmt.Println("Sent permanent move...")
	}
}

func sendTurnEnd(c *websocket.Conn) {
	if PlayerID == -1 {
		fmt.Println("[ERROR] Player ID not yet established. Cannot end turn.")
		return
	}

	turnEndMsg := api.Message{
		Type:    "NEXT_TURN", // You need to update the server to accept this type!
		Payload: nil,
	}

	if err := c.WriteJSON(turnEndMsg); err != nil {
		log.Printf("Error sending PROCESS_NEXT_TURN: %v", err)
	}
	fmt.Println("Sent turn end request...")
}
