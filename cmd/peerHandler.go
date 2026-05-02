package main

import (
	"fmt"
	"net"
)

// PeerService holds dependencies for serving peer TCP connections.
// Add fields here (store, logger, config) as the project grows.
type PeerService struct{}

func (s *PeerService) HandleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		fmt.Print(string(buf[:n]))
	}
}
