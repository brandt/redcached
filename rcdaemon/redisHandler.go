package rcdaemon

import (
	"github.com/brandt/redcached/protocol"
	"gopkg.in/redis.v3"
	"strconv"
)

// `get` handler
//
// Getting multiple keys at the same time:
//
// In Redis, GET is only for getting one key.
// In Memcached, GET is a variadic command, accepting multiple keys.
func RedisGet(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	for _, key := range req.Keys {
		// TODO: Use MGET for multiple keys
		value, err := client.Get(key).Result()
		if err == redis.Nil {
			continue // key did not exist
		} else if err != nil {
			return err
		}
		res.Values = append(res.Values, protocol.McValue{key, "0", []byte(value)})
	}
	res.Response = "END"
	return nil
}

func RedisSet(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	value := req.Value

	// TODO: Also set expiration (currently set to 0)
	err := client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}

	res.Response = "STORED"
	return nil
}

// `add` handler
//
// - Stores the data only if it does not already exist.
// - New items are at the top of the LRU.
// - If an item already exists and an add fails, it promotes the item to the front of the LRU anyway.
func RedisSetNX(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	value := req.Value

	// TODO: Also set expiration (currently set to 0)
	result := client.SetNX(key, value, 0)
	if result.Err() != nil {
		return result.Err()
	}

	if result.Val() {
		res.Response = "STORED"
	} else {
		res.Response = "NOT_STORED"
	}
	return nil
}

func RedisDelete(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	keys := req.Keys
	result := client.Del(keys...)
	if result.Err() != nil {
		return result.Err()
	}
	count := result.Val()

	if count > 0 {
		res.Response = "DELETED"
	} else {
		res.Response = "NOT_FOUND"
	}
	return nil
}

// `incr` handler
//
// Non-existent key behavior:
//
// In Redis, if you INCR a non-existent key, it sets it to zero and then performs the increment.
// In Memcached, it is not valid to increment a key that does not already exist.
//
// Incrementing by arbitrary values:
//
// In Redis, INCR is only for bumping up one. You use INCRBY for more.
// In Memcached, the increment amount is a required argument of INCR.
func RedisIncr(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	increment := req.Increment

	exists := client.Exists(key)
	if !exists.Val() {
		res.Response = "NOT_FOUND"
		return nil
	}

	result := client.IncrBy(key, increment)
	if result.Err() != nil {
		return result.Err()
	}
	val := strconv.FormatInt(result.Val(), 10)

	res.Response = val
	return nil
}

func RedisFlushAll(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	result := client.FlushAll()
	if result.Err() != nil {
		return result.Err()
	}

	res.Response = "OK"
	return nil
}

func RedisVersion(client *redis.Client, req *protocol.McRequest, res *protocol.McResponse) error {
	res.Response = "VERSION redcached-0.1"
	return nil
}

////implement: set/get incr (delete) (flush_all)| stats version
//type RedisHandler struct {
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
//func NewRedisHandler() *RedisHandler {
//	return &RedisHandler{
//		values: make(map[string][]byte),
//	}
//}
//
//
//func (h *RedisHandler) Get(req *protocol.McRequest, res *protocol.McResponse) error {
//	for _, key := range req.Keys {
//		value := h.values[key]
//		// TODO missed
//		res.Values = append(res.Values, protocol.McValue{key, "0", value})
//	}
//	return nil
//}
//
//func (h *RedisHandler) Set(req *protocol.McRequest, res *protocol.McResponse) error {
//	key := req.Key
//	value := req.Value
//	h.values[key] = value
//
//	res.Response = "STORED"
//	return nil
//}
//
//func (h *RedisHandler) Delete(req *protocol.McRequest, res *protocol.McResponse) error {
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
//func (h *RedisHandler) Incr(req *protocol.McRequest, res *protocol.McResponse) error {
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
//func (h *RedisHandler) FlushAll(req *protocol.McRequest, res *protocol.McResponse) error {
//	// TODO clear map
//	res.Response = "OK"
//	return nil
//}
//
//func (h *RedisHandler) Version(req *protocol.McRequest, res *protocol.McResponse) error {
//	res.Response = "VERSION simple-memcached-0.1"
//	return nil
//}
//
//
//func (h *RedisHandler) Stats(req *protocol.McRequest) (*protocol.McResponse, error) {
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
