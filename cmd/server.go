package gokv

import (
	"log"
	"log/slog"
	"net"
)

type Config struct {
	ListenAddr string
	Username   string
	Password   string
}
type Server struct {
	Config
	ln        net.Listener
	msgCh     chan Message
	addPeerCh chan *Peer
	peers     map[*Peer]bool
	quitCh    chan struct{}
}
type Message struct {
	cmd  Command
	peer *Peer
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config:    cfg,
		msgCh:     make(chan Message),
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s.ln = ln
	go s.loop()
	slog.Info("BlinkDB Server Running", "listenAddr", s.ListenAddr)
	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw Message Error", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}
func (s *Server) handleMessage(msg Message) error {
	return nil
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	go peer.readLoop()
}

func (s *Server) Shutdown() {
	close(s.quitCh)
	s.ln.Close()
	for p := range s.peers {
		p.conn.Close()
	}
}
