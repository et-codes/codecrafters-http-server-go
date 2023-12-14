package main

import (
	"github.com/codecrafters-io/http-server-starter-go/logging"
)

const (
	host = "localhost"
	port = "4221"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	s := NewServer(host, port)
	if err := s.Listen(); err != nil {
		logger.Fatal("Error starting server: %v", err)
	}
}
