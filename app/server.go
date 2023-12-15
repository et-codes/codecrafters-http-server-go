package main

import (
	"fmt"
	"io/fs"
	"net"
)

type Server struct {
	host string
	port string
	fs   *fs.FS // Filesystem
}

func NewServer(host, port string, fs *fs.FS) *Server {
	return &Server{host, port, fs}
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error("Failed to bind to port %s", port)
		return err
	}
	logger.Info("Listening on port %s...", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection: ", err.Error())
			return err
		}

		handler := NewHandler(conn)
		go handler.Start()
	}
}
