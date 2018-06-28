package main

import (
	"github.com/sysu-activitypluspc/service-end/ui"
	"os"
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
