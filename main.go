package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Could not open file 'messages.txt': %v", err)
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

// Reads the file "messages.txt" in chunks of 8 bytes, handling lines that may be split across reads. It prints each complete line as it is read.
func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		defer f.Close()

		buffer := make([]byte, 8)
		var line string
		for {
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Unknown read error: %v", err)
				return
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			line += parts[0]
			if len(parts) == 2 {
				lines <- line
				line = parts[1]
			}
		}
		if len(line) > 0 {
			lines <- line
		}
	}()
	return lines
}
