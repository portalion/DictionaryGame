package room

import (
	"encoding/json"
	"net/http"
)

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello from create")
}