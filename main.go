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
	defer file.Close()

	buffer := make([]byte, 8)
	var line string
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unknow read error: %v", err)
		}
		parts := strings.Split(string(buffer[:n]), "\n")
		line += parts[0]
		if len(parts) == 2 {
			fmt.Printf("read: %s\n", line)
			line = parts[1]
		}
	}
	if len(line) > 0 {
		fmt.Printf("read: %s\n", line)
	}
}
