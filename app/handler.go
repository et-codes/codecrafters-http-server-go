package main

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"strings"
)

const (
	appOctStream   = "application/octet-stream"
	CRLF           = "\r\n"
	respOK         = "HTTP/1.1 200 OK\r\n\r\n"
	respNotFound   = "HTTP/1.1 404 Not Found\r\n\r\n"
	statusOK       = "HTTP/1.1 200 OK"
	statusNotFound = "HTTP/1.1 404 Not Found"
	textPlain      = "text/plain"
)

type Handler struct {
	conn net.Conn
	fs   fs.FS
}

func NewHandler(conn net.Conn, fs fs.FS) *Handler {
	return &Handler{conn: conn, fs: fs}
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
	case path[:7] == "/files/":
		filename := path[7:]
		response, err = h.downloadFile(filename)
		if err != nil {
			logger.Error(err.Error())
		}
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

func (h *Handler) downloadFile(filename string) ([]byte, error) {
	if _, err := os.Stat(filename); err != nil {
		return newResponse(statusNotFound, textPlain, ""), fmt.Errorf("could not find file: %v", err)
	}

	b, err := fs.ReadFile(h.fs, filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	return newFileResponse(b), nil
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

func newFileResponse(body []byte) []byte {
	return []byte(fmt.Sprintf(
		"%s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
		statusOK,
		"Content-Type", appOctStream,
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
