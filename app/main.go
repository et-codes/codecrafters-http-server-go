package main

import (
	"os"

	"github.com/et-codes/logging"
)

var logger = logging.New(logging.LevelInfo)

func main() {
	var dir string
	if len(os.Args) >= 3 && os.Args[1] == "--directory" {
		dir = os.Args[2]
		logger.Info("Serving directory %s...", dir)
	}

	s := NewServer("localhost", "4221", dir)
	if err := s.Listen(); err != nil {
		logger.Fatal("Error starting server: %v", err)
	}
}
