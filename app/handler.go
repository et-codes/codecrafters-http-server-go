package main

import "net"

type Handler struct {
	conn net.Conn
}

func NewHandler(conn net.Conn) *Handler {
	return &Handler{conn: conn}
}

func (h *Handler) Start() {
	logger.Info("Handler invoked.")
	defer h.conn.Close()

	response := "HTTP/1.1 200 OK\r\n\r\n"

	_, err := h.conn.Write([]byte(response))
	if err != nil {
		logger.Error("Error writing response: %v", err)
		return
	}
}