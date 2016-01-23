package memcached

import (
	//"fmt"
	//"os"
	"github.com/luxuan/go-memcached-server/protocol"
	"strconv"
)

var dict = make(map[string][]byte)

//implement: set/get incr (delete) (flush_all)| stats version

func DefaultGet(req *protocol.McRequest, res *protocol.McResponse) error {
	for _, key := range req.Keys {
		value := dict[key]
		if len(value) != 0 {
			res.Values = append(res.Values, protocol.McValue{key, "0", value})
		}
	}
	res.Response = "END"
	return nil
}

func DefaultSet(req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	value := req.Value
	dict[key] = value

	res.Response = "STORED"
	return nil
}

func DefaultDelete(req *protocol.McRequest, res *protocol.McResponse) error {
	count := 0
	for _, key := range req.Keys {
		if _, exists := dict[key]; exists {
			delete(dict, key)
			count++
		}
	}
	if count > 0 {
		res.Response = "DELETED"
	} else {
		res.Response = "NOT_FOUND"
	}
	return nil
}

func DefaultIncr(req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	increment := req.Increment
	var base int64
	if value, exists := dict[key]; exists {
		var err error
		base, err = strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return err
		}
	}

	value := strconv.FormatInt(base+increment, 10)
	dict[key] = []byte(value)

	res.Response = value
	return nil
}

func DefaultFlushAll(req *protocol.McRequest, res *protocol.McResponse) error {
	// TODO clear map
	res.Response = "OK"
	return nil
}

func DefaultVersion(req *protocol.McRequest, res *protocol.McResponse) error {
	res.Response = "VERSION simple-memcached-0.1"
	return nil
}

// TODO: Implement 'add' operation
// - Stores the data only if it does not already exist.
// - New items are at the top of the LRU.
// - If an item already exists and an add fails, it promotes the item to the front of the LRU anyway.
//
// This is roughly equivalent with the SETNX operation in Redis.

////implement: set/get incr (delete) (flush_all)| stats version
//type DefaultHandler struct {
//	// TODO lock when goroutine
//	values map[string][]byte
//	/*
//	   TODO do stats in framework, especially for cmd stats
//	   stats   map[string]int
//	   stats:   make(map[string]int),
//	   h.stats["cmd_get"]++
//	   h.stats["get_hits"] += len(res.Values)
//	   h.stats["get_misses"] += len(req.Keys) - len(res.Values)
//	   h.stats["cmd_set"]++
//	*/
//}
//
//func NewDefaultHandler() *DefaultHandler {
//	return &DefaultHandler{
//		values: make(map[string][]byte),
//	}
//}
//
//
//func (h *DefaultHandler) Get(req *protocol.McRequest, res *protocol.McResponse) error {
//	for _, key := range req.Keys {
//		value := h.values[key]
//		// TODO missed
//		res.Values = append(res.Values, protocol.McValue{key, "0", value})
//	}
//	return nil
//}
//
//func (h *DefaultHandler) Set(req *protocol.McRequest, res *protocol.McResponse) error {
//	key := req.Key
//	value := req.Value
//	h.values[key] = value
//
//	res.Response = "STORED"
//	return nil
//}
//
//func (h *DefaultHandler) Delete(req *protocol.McRequest, res *protocol.McResponse) error {
//	count := 0
//	for _, key := range req.Keys {
//		if _, exists := h.values[key]; exists {
//			delete(h.values, key)
//			count++
//		}
//	}
//	if count > 0 {
//		res.Response = "DELETED"
//	} else {
//		res.Response = "NOT_FOUND"
//	}
//	return nil
//}
//
//func (h *DefaultHandler) Incr(req *protocol.McRequest, res *protocol.McResponse) error {
//	key := req.Key
//	increment := req.Increment
//	var base int64
//	if value, exists := h.values[key]; exists {
//		var err error
//		base, err = strconv.ParseInt(string(value), 10, 64)
//		if err != nil {
//			return err
//		}
//	}
//
//	value := strconv.FormatInt(base+increment, 10)
//	h.values[key] = []byte(value)
//
//	res.Response = value
//	return nil
//}
//
//func (h *DefaultHandler) FlushAll(req *protocol.McRequest, res *protocol.McResponse) error {
//	// TODO clear map
//	res.Response = "OK"
//	return nil
//}
//
//func (h *DefaultHandler) Version(req *protocol.McRequest, res *protocol.McResponse) error {
//	res.Response = "VERSION simple-memcached-0.1"
//	return nil
//}
//
//
//func (h *DefaultHandler) Stats(req *protocol.McRequest) (*protocol.McResponse, error) {
//    var b bytes.Buffer
//    b.WriteString("STAT pid ")
//    b.WriteString(strconv.Itoa(os.Getpid()))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT uptime ")
//    b.WriteString(strconv.Itoa(int(time.Now().Sub(startTime).Seconds())))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT cmd_get ")
//    b.WriteString(strconv.Itoa(stats.cmd_get))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT cmd_set ")
//    b.WriteString(strconv.Itoa(stats.cmd_set))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT curr_connections ")
//    b.WriteString(strconv.Itoa(stats.curr_connections))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT total_connections ")
//    b.WriteString(strconv.Itoa(stats.total_connections))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT get_hits ")
//    b.WriteString(strconv.Itoa(stats.get_hits))
//    b.WriteString("\r\n")
//
//    b.WriteString("STAT get_misses ")
//    b.WriteString(strconv.Itoa(stats.get_misses))
//    b.WriteString("\r\n")
//
//    b.WriteString("END")
//
//    return protocol.McResponse{Response: b.String()}, nil
//
//}
