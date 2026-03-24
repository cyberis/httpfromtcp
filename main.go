package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Errorf("Could not open file 'messages.txt': %v", err)
		os.Exit(1)
	}
	defer file.Close()

	buffer := make([]byte, 8)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Errorf("Unknow read error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("read: %s\n", buffer[:n])
	}
}
