package headers

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
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
	// Validate the key to ensure leading or trailing whitespace or illegal characters
	keyLength := len(parts[0])
	if keyLength != len(strings.TrimSpace(parts[0])) {
		return 0, false, errors.New("invalid header line - leading or trailing whitespace in header key")
	}
	if keyLength == 0 {
		return 0, false, errors.New("invalid header line - empty header key")
	}
	for _, r := range parts[0] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune("!#$&'*+-.^_`|~", r) {
			return 0, false, fmt.Errorf("invalid header line - illegal character in header key: %q", r)
		}
	}
	// We should normalize the header key to be case-insensitive, but we should trim whitespace from the value
	parts[0] = strings.ToLower(parts[0])
	parts[1] = strings.TrimSpace(parts[1])
	// If the header key already exists, we should append the new value to the existing value, separated by a comma and a space, as per RFC 7230 Section 3.2.2
	if existingValue, exists := h[parts[0]]; exists {
		parts[1] = existingValue + ", " + parts[1]
	}

	h[parts[0]] = parts[1]

	return len(headerLine) + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	value, exists := h[strings.ToLower(key)]
	return value, exists
}
