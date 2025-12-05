package api

import (
	"encoding/json"
	"log"
	"matetra/engine"
	"net/http"
)

type API struct {
	Game *engine.Game
}

func New(game *engine.Game) *API {
	return &API{Game: game}
}

// starts the server + endpoints
func (a *API) Start() {
	http.HandleFunc("/state", a.handleState)
	http.HandleFunc("/add-player", a.handlePlayer)

	log.Println("API running on :1729")
	log.Fatal(http.ListenAndServe(":1729", nil))
}

// GET /state
func (a *API) handleState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.Game.State)
}

// POST /add-player -> have to implement hashing
func (a *API) handlePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Name string `json:"name"`
		Hash string `json:"hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := a.Game.AddPlayer(body.Name, body.Hash); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "player added",
	})
}
