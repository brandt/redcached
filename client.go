// go-redis-server is a helper library for building server software capable of speaking the redis protocol.
// This could be an alternate implementation of redis, a custom proxy to redis,
// or even a completely different backend capable of "masquerading" its API as a redis database.

package memcached

import (
	"fmt"
    "bufio"
    "strings"
	//"io"
	//"io/ioutil"
	"net"
    "log"
    "time"
	"github.com/luxuan/go-memcached-server/protocol"
)

type HandlerFn func(req *protocol.McRequest) (*protocol.McResponse, error)

type Client struct {
    Addr string     // conn.RemoteAddr().String()
    Conn        net.Conn             // i/o connection
    methods     map[string]HandlerFn // refer to Server.methods
}

// refer to golang/net/http
func NewClient(conn net.Conn, methods map[string]HandlerFn) (c *Client, err error) {
    // TODO set start time

    conn.SetKeepAlive(true)
    conn.SetKeepAlivePeriod(3 * time.Minute)

    return &Client {
        Addr: conn.RemoteAddr().String(),
        Conn: conn,
        methods: methods,
    }, nil
}


// Refer mrproxy/processMc 
func (client *Client) Serve() (err error) {
    conn := client.Conn
	defer func() {
		if err != nil {
			fmt.Fprintf(client.Conn, "-%s\n", err)
		}
		conn.Close()
	}()

    // process
    br := bufio.NewReader(conn)
    bw := bufio.NewWriter(conn)

    for {
        req, err := protocol.ReadRequest(br)
        if perr, ok := err.(protocol.ProtocolError); ok {
            log.Printf("%v ReadRequest protocol err: %v", conn, err)
            bw.WriteString("CLIENT_ERROR " + perr.Error() + "\r\n")
            bw.Flush()
            continue
        } else if err != nil {
            log.Printf("%v ReadRequest err: %v", conn, err)
            return err
        }
        log.Printf("%v Req: %+v\n", conn, req)

        cmd := strings.ToLower(req.Command)
        if cmd == "quit" {
            log.Printf("client send quit, closed")
            return nil
        }

        fn, exists := client.methods[cmd]
        if exists {
            res, err := fn(req)
            if !req.Noreply {
                log.Printf("%v Res: %+v\n", conn, res)
                bw.WriteString(res.Protocol())
                bw.Flush()
            }
        } else {
            res := protocol.McResponse{Response: "ERROR not implement in handler"}
            bw.WriteString(res.Protocol())
            bw.Flush()
        }
    }
	return nil
}
