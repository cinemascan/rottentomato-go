package server

import (
	"flag"

	"github.com/joho/godotenv"
)

// TODO: Add db info when the time comes
type Config struct {
	Port string
}

func NewConfig() *Config {
	godotenv.Load()
	port := flag.String("port", "8081", "port to run the server on; defaults to 8081")
	flag.Parse()
	return &Config{
		Port: *port,
	}
}
