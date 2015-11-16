package memcached

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	DEFAULT_PORT = 11211
)

type Server struct {
	Addr         string // TCP address to listen on, ":6389" if empty
	methods      map[string]HandlerFn
	MonitorChans []chan string

	StartTime        time.Time
	CurrConnections  int
	TotalConnections int
}

// refer from docker/go-redis-server
func NewServer(addr string, methods map[string]HandlerFn) (*Server, error) {
	if addr == "" {
		addr = fmt.Sprintf(":%d", DEFAULT_PORT)
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

	/* //register in handler
	   rh := reflect.TypeOf(handler)
	   for i := 0; i < rh.NumMethod(); i++ {
	       method := rh.Method(i)
	       // inner function
	       if method.Name[0] > 'a' && method.Name[0] < 'z' {
	           continue
	       }
	       // NEED: split Upper
	       fn := &method.Func
	       if fn != nil {
	           key := strings.ToLower(method.Name)
	           log.Printf("REGISTER: %s", key)
	           srv.methods[key] = fn
	       }
	   }
	*/

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
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		client, err := NewClient(conn, srv.methods)
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
