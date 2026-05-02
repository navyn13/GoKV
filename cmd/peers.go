package main

import (
	"errors"
	"log"
	"net"
)

// acceptPeer accepts TCP connections until ln is closed or returns a non-recoverable error.
// Each connection is handled in its own goroutine via handle.
func acceptPeer(ln net.Listener, handle func(net.Conn)) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("accept peer: %v", err)
			continue
		}
		go handle(conn)
	}
}
