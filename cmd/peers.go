package gokv

import (
	"fmt"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
}

// acceptPeer accepts TCP connections until ln is closed or returns a non-recoverable error.
// Each connection is handled in its own goroutine via handle.
func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		data := buf[:n]
		fmt.Println(string(data))
	}
	return fmt.Errorf("connection closed")
}
