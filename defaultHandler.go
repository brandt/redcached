package memcached

import (
	//"fmt"
    //"os"
	"strconv"
	"github.com/luxuan/go-memcached-server/protocol"
)

//implement: set/get incr (delete) (flush_all)| stats version
type DefaultHandler struct {
    // TODO lock when goroutine
	values  map[string][]byte
    /*
    TODO do stats in framework, especially for cmd stats
    stats   map[string]int
    stats:   make(map[string]int),
    h.stats["cmd_get"]++
    h.stats["get_hits"] += len(res.Values)
    h.stats["get_misses"] += len(req.Keys) - len(res.Values)
    h.stats["cmd_set"]++
    */
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{
		values:  make(map[string][]byte),
	}
}

func (h *DefaultHandler) Get(req *protocol.McRequest) (*protocol.McResponse, error) {
    res := &protocol.McResponse{}
    for _, key := range req.Keys {
        value := h.values[key]
        // TODO missed
        res.Values = append(res.Values, protocol.McValue{key, "0", value})
    }
    return res, nil
}

func (h *DefaultHandler) Set(req *protocol.McRequest) (*protocol.McResponse, error) {
    key := req.Key
    value := req.Value
	h.values[key] = value
	return &protocol.McResponse{Response: "STORED"}, nil
}

func (h *DefaultHandler) Delete(req *protocol.McRequest) (*protocol.McResponse, error) {
	count := 0
	for _, key := range req.Keys {
		if _, exists := h.values[key]; exists {
			delete(h.values, key)
			count++
		}
	}
    if count > 0 { 
        return &protocol.McResponse{Response: "DELETED"}, nil
    }
    return &protocol.McResponse{Response: "NOT_FOUND"}, nil
}

func (h *DefaultHandler) Incr(req *protocol.McRequest) (*protocol.McResponse, error) {
    key := req.Key
    increment := req.Increment
    var base int64
    if value, exists := h.values[key]; exists {
        base, err := strconv.ParseInt(string(value), 10, 64) 
        if err != nil {
            return nil, err
        }
    }

    value := strconv.FormatInt(base + increment, 10)
	h.values[key] = []byte(value)
    return &protocol.McResponse{Response: value}, nil
}

func (h *DefaultHandler) FlushAll(req *protocol.McRequest) (*protocol.McResponse, error) {
    return &protocol.McResponse{Response: "OK"}, nil
}

func (h *DefaultHandler) Version(req *protocol.McRequest) (*protocol.McResponse, error) {
    return &protocol.McResponse{Response: "VERSION simple-memcached-0.1"}, nil
}

/*
func (h *DefaultHandler) Stats(req *protocol.McRequest) (*protocol.McResponse, error) {
    var b bytes.Buffer
    b.WriteString("STAT pid ")
    b.WriteString(strconv.Itoa(os.Getpid()))
    b.WriteString("\r\n")

    b.WriteString("STAT uptime ")
    b.WriteString(strconv.Itoa(int(time.Now().Sub(startTime).Seconds())))
    b.WriteString("\r\n")

    b.WriteString("STAT cmd_get ")
    b.WriteString(strconv.Itoa(stats.cmd_get))
    b.WriteString("\r\n")

    b.WriteString("STAT cmd_set ")
    b.WriteString(strconv.Itoa(stats.cmd_set))
    b.WriteString("\r\n")

    b.WriteString("STAT curr_connections ")
    b.WriteString(strconv.Itoa(stats.curr_connections))
    b.WriteString("\r\n")

    b.WriteString("STAT total_connections ")
    b.WriteString(strconv.Itoa(stats.total_connections))
    b.WriteString("\r\n")

    b.WriteString("STAT get_hits ")
    b.WriteString(strconv.Itoa(stats.get_hits))
    b.WriteString("\r\n")

    b.WriteString("STAT get_misses ")
    b.WriteString(strconv.Itoa(stats.get_misses))
    b.WriteString("\r\n")

    b.WriteString("END")

    return protocol.McResponse{Response: b.String()}, nil

}
*/

