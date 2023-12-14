package main

import (
	"bufio"
	"net"
	"strings"
)

const (
	respOK       = "HTTP/1.1 200 OK\r\n\r\n"
	respNotFound = "HTTP/1.1 404 Not Found\r\n\r\n"
)

type Handler struct {
	conn net.Conn
}

func NewHandler(conn net.Conn) *Handler {
	return &Handler{conn: conn}
}

func (h *Handler) Start() {
	logger.Info("Handler invoked.")
	defer h.conn.Close()

	scanner := bufio.NewScanner(h.conn)

	scanner.Scan()
	response := scanner.Text()
	logger.Debug(response)

	parts := strings.Split(response, " ")
	path := parts[1]

	switch path {
	case "/":
		_, err := h.conn.Write([]byte(respOK))
		if err != nil {
			logger.Error("Error writing response: %v", err)
			return
		}
	default:
		_, err := h.conn.Write([]byte(respNotFound))
		if err != nil {
			logger.Error("Error writing response: %v", err)
			return
		}
	}
}
