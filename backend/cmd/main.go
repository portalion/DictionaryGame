package main

import (
	"server/internal/router"
	"server/internal/ws"
)

func main(){
    const host = "localhost"
    const port = 8080

    roomManager := ws.NewRoomManager()

    currentRouter := router.NewRouter()
    currentRouter.SetupMiddleware()
    currentRouter.SetupRoutes(roomManager)
    currentRouter.Start(host, port)
}