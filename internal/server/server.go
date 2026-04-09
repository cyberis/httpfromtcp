package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("could not start TCP listener: %w", err)
	}
	s := &Server{
		listener: l,
	}
	go s.listen()
	return s, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	log.Println("Client connected!")
	response := []byte("HTTP/1.1 200 OK\r\nContent-Length: 12\r\n\r\nHello World!")
	n, err := conn.Write(response)
	if err != nil {
		log.Printf("Could not write response: %v", err)
	} else {
		log.Printf("Wrote %d bytes to client", n)
	}
	log.Println("Client closed!")
}

func (s *Server) Close() error {
	return s.listener.Close()
}
