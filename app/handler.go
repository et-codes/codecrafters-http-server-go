package main

import (
	"fmt"
	"io/fs"
	"net"
	"strings"
)

const (
	appOctStream      = "application/octet-stream"
	CRLF              = "\r\n"
	respOK            = "HTTP/1.1 200 OK\r\n\r\n"
	respNotFound      = "HTTP/1.1 404 Not Found\r\n\r\n"
	respServerError   = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	statusOK          = "HTTP/1.1 200 OK"
	statusNotFound    = "HTTP/1.1 404 Not Found"
	statusServerError = "HTTP/1.1 500 Internal Server Error"
	textPlain         = "text/plain"
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
	if method != "GET" && method != "POST" {
		logger.Warning("Unsupported method %s", method)
	}

	// Generate the response according to the path.
	var response []byte
	switch {
	case path == "/" && method == "GET":
		response = []byte(respOK)
	case path[:6] == "/echo/" && method == "GET":
		message := path[6:]
		response = newTextResponse(statusOK, message)
	case path == "/user-agent" && method == "GET":
		message := parseUserAgent(lines)
		response = newTextResponse(statusOK, message)
	case path[:7] == "/files/":
		filename := path[7:]
		if method == "GET" {
			response, err = h.downloadResponse(filename)
			if err != nil {
				logger.Error(err.Error())
			}
		} else if method == "POST" {
			response, err = h.uploadResponse(filename)
			if err != nil {
				logger.Error(err.Error())
			}
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

func (h *Handler) downloadResponse(filename string) ([]byte, error) {
	logger.Info("User downloading %s...", filename)
	if _, err := fs.Stat(h.fs, filename); err != nil {
		return newTextResponse(statusNotFound, ""), fmt.Errorf("could not find file: %v", err)
	}

	b, err := fs.ReadFile(h.fs, filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	return []byte(fmt.Sprintf(
		"%s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
		statusOK,
		"Content-Type", appOctStream,
		"Content-Length", len(b),
		b,
	)), nil
}

func (h *Handler) uploadResponse(filename string) ([]byte, error) {
	logger.Info("User uploading file %s...", filename)

	// GET THE UPLOAD HERE...

	return nil, nil
}

func newTextResponse(status, body string) []byte {
	return []byte(fmt.Sprintf(
		"%s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
		status,
		"Content-Type", textPlain,
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
