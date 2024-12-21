package main

import (
	"server/router"
)

func main(){
    const host = "localhost"
    const port = 8080
    
    currentRouter := router.NewRouter()
    currentRouter.SetupMiddleware()
    currentRouter.SetupRoutes()
    currentRouter.Start(host, port)
}