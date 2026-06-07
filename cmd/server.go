package gokv

import (
	"errors"
	"log"
	"log/slog"
	"net"
	"strconv"
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
	kv        *KV
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
		quitCh:    make(chan struct{}),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s.ln = ln
	go s.loop()
	slog.Info("GoKV Server Running", "listenAddr", s.ListenAddr)
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
	if _, isAuth := msg.cmd.(AuthCommand); !isAuth {
		if !msg.peer.isAuth {
			msg.peer.conn.Write([]byte("-NOAUTH Authentication required - AUTH {USERNAME} {PASSWORD}\r\n"))
			return nil
		}
	}
	switch v := msg.cmd.(type) {
	case AuthCommand:
		isAuthenticated := s.kv.Auth(v.username, v.password)
		msg.peer.isAuth = isAuthenticated
		var val []byte
		if !isAuthenticated {
			val = []byte("Not Authenticated")
		} else {
			val = []byte("Authenticated")
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	case SetCommand:
		s.kv.Set(v.key, v.val)
	case ExpireCommand:
		s.kv.Expire(v.key, v.seconds)
	case DeleteCommand:
		s.kv.Delete(v.key)
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			val = []byte("key not found")
		}

		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	case ExistCommand:
		isExist, err := s.kv.Exist(v.key)
		if err != nil {
			slog.Error("exist error", "err", err)
		}

		data := []byte(strconv.Itoa(isExist))
		_, err = msg.peer.Send(data)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	case IncrCommand:
		incr, err := s.kv.Incr(v.key)
		if err != nil {
			slog.Error("incr error", "err", err)
		}
		data := []byte(strconv.Itoa(incr))
		_, err = msg.peer.Send(data)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	case DecrCommand:
		decr, err := s.kv.Decr(v.key)
		if err != nil {
			slog.Error("decr error", "err", err)
		}
		data := []byte(strconv.Itoa(decr))
		_, err = msg.peer.Send(data)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	}

	return nil
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
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
