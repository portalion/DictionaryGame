package router

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/ws"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {return true},
	}
type Router struct{
	router *mux.Router

	hub *ws.RoomHub
}

func NewRouter(hub *ws.RoomHub) *Router{
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