package main

import (
	"os"

	"github.com/sysu-activitypluspc/service-end/ui"
)

var PORT = "8080"

func main() {
	port := os.Getenv("PORT")
	if len(port) != 0 {
		PORT = port
	}
	PORT = ":" + PORT

	server := ui.GetServer()
	server.Run(PORT)
}
