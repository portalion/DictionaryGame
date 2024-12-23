package router

import (
	"fmt"
	"log"
	"net/http"
	"server/ws"

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

func joinRoomHandler(rm *ws.Room, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	rm.JoinRoom(username, conn)
}

func (r *Router) SetupRoutes(rm *ws.Room) {
	//r.router.HandleFunc("/room/create", room.CreateRoomHandler).Methods(http.MethodGet)
	r.router.HandleFunc("/ws/room/join/{id}", func(w http.ResponseWriter, r *http.Request) {joinRoomHandler(rm, w, r)})
}

func (r *Router) Start(hostname string, port int) {
	url := fmt.Sprintf("%s:%d", hostname, port)
	fmt.Printf("Starting server at: %s\n", url)
	log.Fatal(http.ListenAndServe(url, r.router))
}