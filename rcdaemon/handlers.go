package rcdaemon

import (
	"github.com/brandt/redcached/protocol"
	"gopkg.in/redis.v3"
	"strconv"
	"time"
)

var backend *redis.Client

func init() {
	backend = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		PoolSize: 100,
	})
}

func expirationParser(t uint32) (time.Duration, error) {
	if t > 2592000 { // maximum non-epoch value is 30 days
		now := time.Now()
		expire_at := time.Unix(int64(t), 0)
		return expire_at.Sub(now), nil
	} else {
		secs := time.Duration(t) * time.Second
		return secs, nil
	}
}

// `get` handler
//
// Getting multiple keys at the same time:
//
// In Redis, GET is only for getting one key.
// In Memcached, GET is a variadic command, accepting multiple keys.
func GetHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	for _, key := range req.Keys {
		// TODO: Use MGET for multiple keys
		value, err := backend.Get(key).Result()
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

func SetHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	value := req.Value

	// TODO: Also set expiration (currently set to 0)
	err := backend.Set(key, value, 0).Err()
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
func AddHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	value := req.Value

	// TODO: Also set expiration (currently set to 0)
	result := backend.SetNX(key, value, 0)
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

func DeleteHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	keys := req.Keys
	result := backend.Del(keys...)
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
func IncrHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	key := req.Key
	increment := req.Increment

	exists := backend.Exists(key)
	if !exists.Val() {
		res.Response = "NOT_FOUND"
		return nil
	}

	result := backend.IncrBy(key, increment)
	if result.Err() != nil {
		return result.Err()
	}
	val := strconv.FormatInt(result.Val(), 10)

	res.Response = val
	return nil
}

func FlushAllHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	result := backend.FlushAll()
	if result.Err() != nil {
		return result.Err()
	}

	res.Response = "OK"
	return nil
}

func VersionHandler(req *protocol.McRequest, res *protocol.McResponse) error {
	res.Response = "VERSION redcached-0.1"
	return nil
}
