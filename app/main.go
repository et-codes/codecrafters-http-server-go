package main

import (
	"os"

	"github.com/codecrafters-io/http-server-starter-go/logging"
)

const (
	host          = "localhost"
	port          = "4221"
	directoryFlag = "--directory"
)

var logger = logging.New(logging.LevelDebug)

func main() {
	var dir string
	if len(os.Args) >= 3 && os.Args[1] == directoryFlag {
		dir = os.Args[2]
		logger.Info("Serving directory %s...", dir)
	}

	s := NewServer(host, port, dir)
	if err := s.Listen(); err != nil {
		logger.Fatal("Error starting server: %v", err)
	}
}
