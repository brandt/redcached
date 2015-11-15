// go-redis-server is a helper library for building server software capable of speaking the redis protocol.
// This could be an alternate implementation of redis, a custom proxy to redis,
// or even a completely different backend capable of "masquerading" its API as a redis database.

package memcached

import (
	"fmt"
	//"io"
	//"io/ioutil"
    "log"
    "strings"
    "time"
	"net"
    "reflect"
)

const (
    DEFAULT_PORT = 11211
)

type Server struct {
	Addr         string // TCP address to listen on, ":6389" if empty
    methods      map[string]HandlerFn
	MonitorChans []chan string

    StartTime    time.Time
    CurrConnections  int
    TotalConnections int
}

// refer from docker/go-redis-server
func NewServer(addr string, handler interface{}) (*Server, error) {
    if addr == "" {
        addr = fmt.Sprintf(":%d", DEFAULT_PORT)
    }
	if handler == nil {
		handler = NewDefaultHandler()
	}

	srv := &Server{
        Addr: addr,
        methods:      make(map[string]HandlerFn),
		MonitorChans: []chan string{},

        StartTime: time.Now(),
        CurrConnections: 0,
        TotalConnections: 0,
	}

    /* TODO register in handler
    rh := reflect.TypeOf(handler)
    for i := 0; i < rh.NumMethod(); i++ {
        method := rh.Method(i)
        // inner function
        if method.Name[0] > 'a' && method.Name[0] < 'z' {
            continue
        }
        // TODO split Upper
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
	return srv.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	srv.MonitorChans = []chan string{}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
        client := NewClient(conn, srv.methods)
		go client.Serve()
	}
    return nil
}

