package rcdaemon

import (
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"net"
	"time"
)

const (
	DEFAULT_PORT = 11212
)

type Server struct {
	Addr         string // TCP address to listen on, ":11212" if empty
	methods      map[string]HandlerFn
	MonitorChans []chan string

	StartTime        time.Time
	CurrConnections  int
	TotalConnections int
}

func NewServer(addr string, methods map[string]HandlerFn) (*Server, error) {
	if addr == "" {
		addr = fmt.Sprintf("127.0.0.1:%d", DEFAULT_PORT)
	}
	if methods == nil {
		methods = make(map[string]HandlerFn)
	}

	srv := &Server{
		Addr:         addr,
		methods:      methods,
		MonitorChans: []chan string{},

		StartTime:        time.Now(),
		CurrConnections:  0,
		TotalConnections: 0,
	}

	return srv, nil
}

func (srv *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	log.Printf("Start and Listening at %s", srv.Addr)
	return srv.Serve(l)
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	srv.MonitorChans = []chan string{}
	backend := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		PoolSize: 100,
	})
	defer backend.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		client, err := NewClient(backend, conn, srv)
		if err != nil {
			log.Printf("New Client ERROR:: %v", err)
			continue
		}
		log.Printf("Client %s Connected", client.Addr)
		go client.Serve()
	}
	return nil
}

func (srv *Server) RegisterFunc(name string, fn HandlerFn) error {
	log.Printf("REGISTER func: %s", name)
	srv.methods[name] = fn
	return nil
}
