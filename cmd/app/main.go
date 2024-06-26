package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/dGilli/gokv/pkg/client"
	"github.com/dGilli/gokv/pkg/proto"
)

const defaultListenAddr = ":5301"

type (
	Config struct {
		ListenAddr string
	}

	Server struct {
		Config
		peers     map[*Peer]bool
		ln        net.Listener
		addPeerCh chan *Peer
        quitCh    chan struct{}
        msgCh     chan []byte
	}
)

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
        msgCh:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

    go s.loop()

    slog.Info("server running", "listenAddr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) handleRawMessage(rawMsg []byte) error {
    cmd, err := proto.ParseCommand(string(rawMsg))
    if err != nil {
        return err
    }
    switch v := cmd.(type) {
    case proto.SetCommand:
        slog.Info("wanna set a key in to the hash table", "key", v)
    }

    return nil
}

func (s *Server) loop() {
	for {
		select {
        case rawMsg := <- s.msgCh:
            if err := s.handleRawMessage(rawMsg); err != nil {
                slog.Error("raw message error", "err", err)
            }
		case <-s.quitCh:
            return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
    peer := NewPeer(conn, s.msgCh)
    s.addPeerCh <- peer
    slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
    if err := peer.readLoop(); err != nil {
        slog.Error("peer read error", err, "remoteAddr", conn.RemoteAddr())
    }
}

func main() {
    go func() {
        server := NewServer(Config{})
        log.Fatal(server.Start())
    }()
    time.Sleep(time.Second)

    for i := 0; i < 10; i++ {
        client := client.New("localhost:5301")
        if err := client.Set(context.Background(), "foo", "bar"); err != nil {
            log.Fatal(err)
        }
    }

    time.Sleep(time.Second)
}
