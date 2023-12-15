package main

import (
	"fmt"
	"net"
)

type Server struct {
	host string
	port string
}

func NewServer(host, port string) *Server {
	return &Server{host, port}
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error("Failed to bind to port %s", port)
		return err
	}
	logger.Info("Listening on port %s", port)

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
