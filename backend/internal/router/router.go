package router

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/ws/room"

	"github.com/gorilla/mux"
)
type Router struct{
	router *mux.Router

	hub *room.RoomHub
}

func NewRouter(hub *room.RoomHub) *Router{
	muxRouter := mux.NewRouter()
	muxRouter.Use()

	return &Router{
		router: muxRouter,
		hub: hub,
	}
}

func (r *Router) SetupMiddleware() {
	r.router.Use(corsMiddleware)
	r.router.Use(jsonContentMiddleware)
}

func (r *Router) SetupRoutes() {
	r.router.HandleFunc("/room/create", r.createRoomHandler).Methods(http.MethodPost)
	r.router.HandleFunc("/ws/room/{code}/join", r.joinRoomHandler)
}

func (r *Router) Start(hostname string, port int) {
	url := fmt.Sprintf("%s:%d", hostname, port)
	fmt.Printf("Starting server at: %s\n", url)
	log.Fatal(http.ListenAndServe(url, r.router))
}