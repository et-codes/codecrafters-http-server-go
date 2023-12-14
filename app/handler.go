package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	CRLF           = "\r\n"
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

	// Wait for client request.
	request, err := h.getRequest()
	if err != nil {
		logger.Error("Error reading response: %v", err)
		return
	}

	// Get start line method and path.
	lines := strings.Split(request, CRLF)
	startLine := lines[0]
	logger.Info(startLine)

	words := strings.Split(startLine, " ")
	method, path := words[0], words[1]

	// Make sure we can handle the request type.
	if method != "GET" {
		logger.Warning("Unsupported method %s", method)
	}

	// Generate the response according to the path.
	var response []byte
	switch {
	case path == "/":
		response = []byte(respOK)
	case path[:6] == "/echo/":
		message := path[6:]
		response = newResponse(statusOK, textPlain, message)
	case path == "/user-agent":
		message := parseUserAgent(lines)
		response = newResponse(statusOK, textPlain, message)
	default:
		response = []byte(respNotFound)
	}

	// Send the response.
	_, err = h.conn.Write(response)
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

func parseUserAgent(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "User-Agent: ") {
			return strings.TrimPrefix(line, "User-Agent: ")
		}
	}
	return ""
}
