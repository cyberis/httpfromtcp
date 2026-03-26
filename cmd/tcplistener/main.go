package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {

	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Could not start TCP listener: %v", err)
	}
	defer l.Close()
	fmt.Println("Server is listening on port 42069...")

	// Wait for connects and then read from them
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}

		// Handle the connection in a new goroutine to received write text to stdout
		go func(c net.Conn) {
			// Handle the connection here
			fmt.Println("Client connected!")
			for line := range getLinesChannel(c) {
				fmt.Printf("%s", line)
			}
			fmt.Println("")
			c.Close()
			fmt.Println("Client closed!")
		}(conn)
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
