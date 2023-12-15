package main

import (
	"io/fs"
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
	var fs fs.FS
	if len(os.Args) >= 3 && os.Args[1] == directoryFlag {
		fs = os.DirFS(os.Args[2])
		logger.Info("Serving directory %s...", fs)
	}

	s := NewServer(host, port, &fs)
	if err := s.Listen(); err != nil {
		logger.Fatal("Error starting server: %v", err)
	}
}
