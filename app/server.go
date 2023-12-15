package main

import (
	"fmt"
	"net"
)

type Server struct {
	host      string
	port      string
	directory string
}

func NewServer(host, port, dir string) *Server {
	return &Server{host, port, dir}
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		logger.Error("Failed to bind to port %s", s.port)
		return err
	}
	logger.Info("Listening on port %s...", s.port)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection: ", err.Error())
			return err
		}

		handler := NewHandler(conn, s.directory)
		go handler.Start()
	}
}
