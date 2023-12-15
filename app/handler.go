package main

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	appOctStream      = "application/octet-stream"
	textPlain         = "text/plain"
	CRLF              = "\r\n"
	statusOK          = "HTTP/1.1 200 OK"
	statusCreated     = "HTTP/1.1 201 Created"
	statusNotFound    = "HTTP/1.1 404 Not Found"
	statusMethod      = "HTTP/1.1 405 Method Not Allowed"
	statusServerError = "HTTP/1.1 500 Internal Server Error"
	respOK            = statusOK + CRLF + CRLF
	respCreated       = statusCreated + CRLF + CRLF
	respNotFound      = statusNotFound + CRLF + CRLF
	respMethod        = statusMethod + CRLF + CRLF
	respServerError   = statusServerError + CRLF + CRLF
)

type Handler struct {
	conn      net.Conn
	directory string
}

type Request struct {
	StartLine string
	Method    string
	Path      string
	Header    []string
	Body      []byte
}

func NewHandler(conn net.Conn, dir string) *Handler {
	return &Handler{conn: conn, directory: dir}
}

func (h *Handler) Start() {
	logger.Info("Handler invoked.")
	defer h.conn.Close()

	// Wait for client req.
	req, err := h.getRequest()
	if err != nil {
		logger.Error("Error reading response: %v", err)
		return
	}

	// Generate the response according to the path.
	var response []byte
	switch {
	case req.Method != "GET" && req.Method != "POST":
		response = []byte(respMethod)
	case req.Path == "/" && req.Method == "GET":
		response = []byte(respOK)
	case req.Path[:6] == "/echo/" && req.Method == "GET":
		message := req.Path[6:]
		response = newTextResponse(statusOK, message)
	case req.Path == "/user-agent" && req.Method == "GET":
		message := parseUserAgent(req.Header)
		response = newTextResponse(statusOK, message)
	case req.Path[:7] == "/files/":
		filename := req.Path[7:]
		if req.Method == "GET" {
			response, err = h.downloadResponse(filename)
			if err != nil {
				logger.Error(err.Error())
			}
		} else if req.Method == "POST" {
			response, err = h.uploadResponse(filename, req.Body)
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

func (h *Handler) getRequest() (*Request, error) {
	request := make([]byte, 16*1024) // 16kB
	n, err := h.conn.Read(request)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Get start line method and path.
	lines := strings.Split(string(request[:n]), CRLF)
	startLine := lines[0]
	logger.Info(startLine)

	words := strings.Split(startLine, " ")
	method, path := words[0], words[1]

	// Get message body.
	var body []byte
	var header []string
	for i, line := range lines {
		if line == "" {
			body = []byte(lines[i+1])
			break
		}
		header = append(header, line)
	}

	return &Request{
		StartLine: startLine,
		Method:    method,
		Path:      path,
		Header:    header,
		Body:      body,
	}, nil
}

// downloadResponse returns a 200 OK response message with a header and the
// contents of the specified file in the body, or 404 if the file is not found.
func (h *Handler) downloadResponse(filename string) ([]byte, error) {
	if h.directory == "" {
		return []byte(respNotFound), fmt.Errorf("directory not specified")
	}

	// Open directory.
	fileSystem := os.DirFS(h.directory)
	if _, err := fs.Stat(fileSystem, filename); err != nil {
		return []byte(respNotFound), fmt.Errorf("could not find file: %v", err)
	}

	// Read the file.
	logger.Info("User downloading %s...", filename)
	b, err := fs.ReadFile(fileSystem, filename)
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

// uploadResponse receives a file from the client and saves it, returning
// a 201 Created response if successful, 500 if not.
func (h *Handler) uploadResponse(filename string, data []byte) ([]byte, error) {
	// Create the file.
	path := filepath.Join(h.directory, filename)
	file, err := os.Create(path)
	if err != nil {
		return []byte(respServerError),
			fmt.Errorf("error creating file %s: %v", path, err)
	}
	defer file.Close()

	// Write to the file.
	logger.Info("User uploading file %s...", filename)
	if _, err = file.Write(data); err != nil {
		return []byte(respServerError),
			fmt.Errorf("error writing file %s: %v", path, err)
	}

	return []byte(respCreated), nil
}

// newTestResponse returns an HTTP response with a text body.
func newTextResponse(status, body string) []byte {
	return []byte(fmt.Sprintf(
		"%s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
		status,
		"Content-Type", textPlain,
		"Content-Length", len(body),
		body,
	))
}

// parseUserAgent extracts the User-Agent value from the headers.
func parseUserAgent(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "User-Agent: ") {
			return strings.TrimPrefix(line, "User-Agent: ")
		}
	}
	return ""
}
