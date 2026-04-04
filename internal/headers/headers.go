package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

const (
	CRLF = "\r\n"
)

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// We need to read until we have a complete header section, which is terminated by two CRLFs in a row
	headerSection := string(data)
	if !strings.Contains(headerSection, CRLF) {
		return 0, false, nil // Not enough data to parse headers yet
	}
	if strings.HasPrefix(headerSection, CRLF) {
		return 2, true, nil // No headers, just the CRLFs separating headers and body
	}
	headerLine := strings.Split(headerSection, CRLF)[0]
	parts := strings.SplitN(headerLine, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header line - missing colon")
	}
	keyLength := len(parts[0])
	if keyLength != len(strings.TrimSpace(parts[0])) {
		return 0, false, errors.New("invalid header line - leading or trailing whitespace in header key")
	}
	parts[1] = strings.TrimSpace(parts[1])

	h[parts[0]] = parts[1]

	return len(headerLine) + 2, false, nil
}
