package main

import (
	"server/internal/router"
	"server/internal/ws/room"
)

func main(){
    const host = "localhost"
    const port = 8080

    hub := room.NewRoomHub()

    currentRouter := router.NewRouter(hub)
    currentRouter.SetupMiddleware()
    currentRouter.SetupRoutes()
    currentRouter.Start(host, port)
}