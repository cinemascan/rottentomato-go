package main

import (
	"github.com/cinemascan/rottentomato-server/internal/app/server"
)

func main() {
	config := server.NewConfig()
	server.Init(config.Port)
}
