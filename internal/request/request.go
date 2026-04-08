package request

import (
	"errors"
	"io"
	"strings"

	"github.com/cyberis/httpfromtcp/internal/headers"
)

// We need to create an enum to track parser state
type parserState int

const (
	initialized parserState = iota
	parseHeaders
	done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	ParserState parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	CRLF         = "\r\n"
	bufferSize   = 8
	validmethods = "GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH|TRACE|CONNECT"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{
		ParserState: initialized,
		Headers:     headers.NewHeaders(),
	}

	buf := make([]byte, bufferSize) // Start with a small buffer to read the request line
	readToIndex := 0

	// We need to read until we have a complete request line, which is terminated by CRLF
	for {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF // Not enough data to parse the request line
			}
			return nil, err
		}
		if n == 0 {
			return nil, io.ErrUnexpectedEOF // No data read, but not EOF
		}
		readToIndex += n
		bytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if bytesParsed > 0 {
			copy(buf, buf[bytesParsed:readToIndex])
			readToIndex -= bytesParsed
		}
		if request.ParserState == done {
			break
		}
	}

	return &request, nil
}

func parseRequestLine(data string) (*RequestLine, int, error) {
	var requestLine RequestLine

	// Do we have the whole request line? It should include a CRLF at the end, but we can also check for the presence of three parts (method, target, version)
	if !strings.HasSuffix(data, CRLF) {
		return nil, 0, nil // Not enough data to parse the request line yet
	}
	line := strings.Split(data, CRLF)
	parts := strings.Split(line[0], " ")
	if len(parts) != 3 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	// Validate the Method as one of the standard HTTP methods
	requestLine.Method = parts[0]
	if !strings.Contains(validmethods, requestLine.Method) {
		return nil, 0, errors.New("invalid HTTP method")
	}
	// Validate the Request Target (basic validation, can be improved)
	requestLine.RequestTarget = parts[1]
	if requestLine.RequestTarget == "" {
		return nil, 0, errors.New("invalid request target")
	}
	// Validate the HTTP Version, which currently must be 1.1 only
	requestLine.HttpVersion = parts[2]
	if !strings.HasPrefix(requestLine.HttpVersion, "HTTP/") {
		return nil, 0, errors.New("invalid HTTP version")
	}
	requestLine.HttpVersion = strings.TrimPrefix(requestLine.HttpVersion, "HTTP/")
	if requestLine.HttpVersion != "1.1" {
		return nil, 0, errors.New("unsupported HTTP version")
	}

	return &requestLine, len(line[0]) + len(CRLF), nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.ParserState != done {
		bytesParsed, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += bytesParsed
		if bytesParsed == 0 {
			break
		}
	}
	return totalBytesParsed, nil

}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.ParserState {
	case initialized:
		requestLine, bytesParsed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if requestLine == nil {
			return 0, nil // Not enough data to parse the request line yet
		}
		r.RequestLine = *requestLine
		r.ParserState = parseHeaders
		return bytesParsed, nil
	case parseHeaders:
		bytesParsed, headersDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if headersDone {
			r.ParserState = done
		}
		return bytesParsed, nil
	case done:
		return 0, errors.New("request already parsed")
	default:
		return 0, errors.New("invalid parser state")
	}

}
