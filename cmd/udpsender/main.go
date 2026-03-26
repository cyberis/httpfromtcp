package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Could not resolve UDP address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Could not dial UDP: %v", err)
	}
	defer conn.Close()
	fmt.Println("UDP sender is ready. Type messages to send to the UDP listener (Ctrl+C to exit).")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		_, err = conn.Write([]byte(text))
		if err != nil {
			log.Printf("Error sending UDP message: %v", err)
			continue
		}
		fmt.Printf("Message Sent: %s", text)
	}
}
