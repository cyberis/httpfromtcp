package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cyberis/httpfromtcp/internal/request"
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
			request, err := request.RequestFromReader(c)
			if err != nil {
				log.Printf("Could not parse request: %v", err)
			} else {
				fmt.Println("Request line:")
				fmt.Printf("- Method: %s\n", request.RequestLine.Method)
				fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
				fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
			}
			c.Close()
			fmt.Println("Client closed!")
		}(conn)
	}
}
