package router

import (
	"fmt"
	"log"
	"net/http"
	"server/room"

	"github.com/gorilla/mux"
)

type Router struct{
	router *mux.Router
}

func NewRouter() *Router{
	muxRouter := mux.NewRouter()
	muxRouter.Use()

	return &Router{
		router: muxRouter,
	}
}

func (r *Router) SetupMiddleware() {
	r.router.Use(corsMiddleware)
	r.router.Use(jsonContentMiddleware)
}

func (r *Router) SetupRoutes() {
	r.router.HandleFunc("/room/create", room.CreateRoomHandler).Methods(http.MethodGet)
	r.router.HandleFunc("/room/join/{id}", nil)
}

func (r *Router) Start(hostname string, port int) {
	url := fmt.Sprintf("%s:%d", hostname, port)
	fmt.Printf("Starting server at: %s\n", url)
	log.Fatal(http.ListenAndServe(url, r.router))
}