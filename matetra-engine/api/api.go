package api

import (
	"encoding/json"
	"log"
	"matetra/engine"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PlayerPayload struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type CardPlayPayload struct {
	CardIndex int   `json:"card_index"`
	Inputs    []int `json:"inputs"`
	Permanent bool  `json:"permanent"`
}

type PlayerConnection struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	PlayerID int
}

type API struct {
	Game        *engine.Game
	Connections map[int]*PlayerConnection
	nextConnID  int
}

func New(game *engine.Game) *API {
	return &API{
		Game:        game,
		Connections: make(map[int]*PlayerConnection),
	}
}

// websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// starts the server + endpoints
func (a *API) Start() {
	http.HandleFunc("/ws", a.handleWebSocket)

	log.Println("API running on :1729")
	log.Fatal(http.ListenAndServe(":1729", nil))
}

func (a *API) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade connection: %v", err)
		return
	}

	playerConn := &PlayerConnection{conn: conn, PlayerID: -1}
	connID := a.nextConnID
	a.Connections[connID] = playerConn
	a.nextConnID++

	log.Printf("client %d connected", connID)

	go a.readMessages(connID, playerConn)
}

func (a *API) readMessages(connID int, pc *PlayerConnection) {
	defer func() {
		pc.conn.Close()
		log.Printf("client %d disconnected", connID)
	}()

	for {
		var incomingMsg Message
		if err := pc.conn.ReadJSON(&incomingMsg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				return
			}
			log.Printf("read error for client %d, %v", connID, err)
			return
		}
		a.handleIncomingMessages(pc, incomingMsg)
	}
}

func (a *API) handleIncomingMessages(pc *PlayerConnection, msg Message) {
	switch msg.Type {
	case "ADD_PLAYER":
		var payload PlayerPayload
		payloadBytes, err := json.Marshal(msg.Payload)
		if err != nil {
			log.Printf("error marshalling payload: %v", err)
		}
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			a.sendError(pc, "invalid player payload format")
			return
		}

		playerID := len(a.Game.State.Players)

		if err := a.Game.AddPlayer(payload.Name, payload.Hash); err != nil {
			a.sendError(pc, err.Error())
			return
		}

		pc.PlayerID = playerID

		a.sendResponse(pc, "PLAYER_ADDED", map[string]string{"name": payload.Name})
		a.BroadcastState()
	case "PLAY_CARD":
		a.handlePlayCard(pc, msg.Payload)
	default:
		a.sendError(pc, "unknown message type: "+msg.Type)
	}
}

func (a *API) handlePlayCard(pc *PlayerConnection, payload interface{}) {
	if pc.PlayerID == -1 {
		a.sendError(pc, "player is not authenticated")
		return
	}

	var cardPayload CardPlayPayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		a.sendError(pc, "error parsing card play payload")
		return
	}
	if err := json.Unmarshal(payloadBytes, &cardPayload); err != nil {
		a.sendError(pc, "invalid card play payload format")
		return
	}

	// Card ownership check
	if !a.Game.PlayerCanPlayCard(pc.PlayerID, cardPayload.CardIndex) {
		a.sendError(pc, "you do not own this card")
		return
	}

	if cardPayload.Permanent {
		// We
	} else {

	}
}

func (a *API) BroadcastState() {
	stateMsg := Message{
		Type:    "STATE_UPDATE",
		Payload: a.Game.State,
	}
	for _, pc := range a.Connections {
		pc.mu.Lock()
		if err := pc.conn.WriteJSON(stateMsg); err != nil {
			log.Printf("error broadcasting state: %v", err)
		}
		pc.mu.Unlock()
	}
}

func (a *API) sendResponse(pc *PlayerConnection, responseType string, data interface{}) {
	respMsg := Message{
		Type:    responseType,
		Payload: data,
	}
	pc.mu.Lock()
	if err := pc.conn.WriteJSON(respMsg); err != nil {
		log.Printf("error sending response: %v", err)
	}
	pc.mu.Unlock()
}

func (a *API) sendError(pc *PlayerConnection, errMsg string) {
	errorMsg := Message{
		Type: "ERROR",
		Payload: map[string]string{
			"message": errMsg,
		},
	}
	pc.mu.Lock()
	if err := pc.conn.WriteJSON(errorMsg); err != nil {
		log.Printf("error sending error: %v", err)
	}
	pc.mu.Unlock()
}
