package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on :8080")
	svc := &PeerService{}
	acceptPeer(ln, svc.HandleConn)
}
