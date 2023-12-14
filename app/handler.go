package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	CRLF           = "\r\n"
	delim          = "\r\n\r\n"
	respOK         = "HTTP/1.1 200 OK\r\n\r\n"
	respNotFound   = "HTTP/1.1 404 Not Found\r\n\r\n"
	statusOK       = "HTTP/1.1 200 OK"
	statusNotFound = "HTTP/1.1 404 Not Found"
	textPlain      = "text/plain"
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

	// Get
	request, err := h.getRequest()
	if err != nil {
		logger.Error("Error reading response: %v", err)
		return
	}

	lines := strings.Split(request, CRLF)
	startLine := lines[0]
	logger.Info(startLine)

	words := strings.Split(startLine, " ")
	method, path := words[0], words[1]

	if method != "GET" {
		logger.Warning("Unsupported method %s", method)
	}

	var response []byte
	switch {
	case path == "/":
		response = []byte(respOK)
		_, err = h.conn.Write(response)
	case path[:6] == "/echo/":
		message := path[6:]
		response = newResponse(statusOK, textPlain, message)
		_, err = h.conn.Write(response)
	default:
		response = []byte(respNotFound)
		_, err = h.conn.Write(response)
	}
	if err != nil {
		logger.Error("Error writing reply: %v", err)
		return
	}
	logger.Debug("reply sent.\n%s", response)
}

func (h *Handler) getRequest() (string, error) {
	response := make([]byte, 1024)
	_, err := h.conn.Read(response)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	return string(response), nil
}

func newResponse(status, contentType, body string) []byte {
	return []byte(fmt.Sprintf(
		"%s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
		status,
		"Content-Type", contentType,
		"Content-Length", len(body),
		body,
	))
}
