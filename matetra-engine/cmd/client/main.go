package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"matetra/model"
	"matetra/utils"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)

	cmd := os.Args

	if len(cmd) < 2 {
		fmt.Println("usage: go run matetra <server-address>")
		fmt.Println("example: go run matetra http://localhost:1729")
		return
	}

	serverAddr := cmd[1]

	if !strings.HasPrefix(serverAddr, "http") {
		serverAddr = "http://" + serverAddr
	}

	fmt.Printf("attempting to connect to server at %s...\n", serverAddr)
	if !checkServerStatus(serverAddr) {
		log.Fatalf("error: could not connect to server at %s. Is the server running?", serverAddr)
	}
	fmt.Println("Connection successful.")

	addPlayer(serverAddr)
}

func checkServerStatus(addr string) bool {
	resp, err := http.Get(addr + "/state")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func addPlayer(addr string) {
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
		log.Fatal("username and password cannot be empty.")
	}

	passwordHash := utils.Hash(password)
	fmt.Printf("hashed password (SHA256): %s...\n", passwordHash[:8])

	// construct the Player object and JSON payload
	playerData := model.Player{
		Name: username,
		Hash: passwordHash,
	}

	jsonPayload, err := json.Marshal(playerData)
	if err != nil {
		log.Fatalf("error creating JSON payload: %v", err)
	}

	// send POST request to /add-player
	fmt.Println("registering player...")
	resp, err := http.Post(addr+"/add-player", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("error sending registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("registration failed with status code %d\nserver eror: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	fmt.Printf("success: player @%s has beed added to the game!\n", username)
}
