package main

import (
	"server/router"
	"server/ws"
)

func main(){
    const host = "localhost"
    const port = 8080

    roomManager := ws.NewRoom()
    
    currentRouter := router.NewRouter()
    currentRouter.SetupMiddleware()
    currentRouter.SetupRoutes(roomManager)
    currentRouter.Start(host, port)
}