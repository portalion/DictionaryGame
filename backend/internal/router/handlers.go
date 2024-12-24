package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (router *Router) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := router.hub.CreateRoom()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(CreateRoomResponse{Code: roomCode}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (router *Router) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		log.Println("Username is required")
		return
	}

	code := mux.Vars(r)["code"]
	if code == "" {
		http.Error(w, "Room code is required", http.StatusBadRequest)
		log.Println("Room code is required")
		return
	}

	err := router.hub.JoinRoom(w, r, code, username)

	if err != nil {
		http.Error(w, "Failed to join the room", http.StatusNotFound)
		log.Println("Failed to join the room")
		return
	}
}
