package main

import (
	"os"
	"github.com/sysu-activitypluspc/service-end/router"
)

var PORT = "8080"

func main() {
	port := os.Getenv("PORT")
	if len(port) != 0 {
		PORT = port
	}
	PORT = ":" + PORT

	server := router.GetServer()
	server.Run(PORT)
}