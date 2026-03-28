package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	CRLF         = "\r\n"
	validmethods = "GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH|TRACE|CONNECT"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request

	requestBuffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(requestBuffer) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	lines := strings.Split(string(requestBuffer), CRLF)
	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}
	request.RequestLine = *requestLine

	return &request, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	var requestLine RequestLine
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, io.ErrUnexpectedEOF
	}
	// Validate the Method as one of the standard HTTP methods
	requestLine.Method = parts[0]
	if !strings.Contains(validmethods, requestLine.Method) {
		return nil, errors.New("invalid HTTP method")
	}
	// Validate the Request Target (basic validation, can be improved)
	requestLine.RequestTarget = parts[1]
	if requestLine.RequestTarget == "" {
		return nil, errors.New("invalid request target")
	}
	// Validate the HTTP Version, which currently must be 1.1 only
	requestLine.HttpVersion = parts[2]
	if !strings.HasPrefix(requestLine.HttpVersion, "HTTP/") {
		return nil, errors.New("invalid HTTP version")
	}
	requestLine.HttpVersion = strings.TrimPrefix(requestLine.HttpVersion, "HTTP/")
	if requestLine.HttpVersion != "1.1" {
		return nil, errors.New("unsupported HTTP version")
	}

	return &requestLine, nil
}
