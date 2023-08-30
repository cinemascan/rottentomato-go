package main

import (
	"github.com/cinemascan/rottentomato-go/internal/app/server"
)

func main() {
	config := server.NewConfig()
	server.Init(config.Port)
}
