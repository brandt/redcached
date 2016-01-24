package main

import (
	"github.com/luxuan/go-memcached-server"
)

func main() {
	server, err := memcached.NewServer("", nil)
	if err != nil {
		panic(err)
	}

	// register handler
	server.RegisterFunc("get", memcached.RedisGet)
	server.RegisterFunc("set", memcached.RedisSet)
	server.RegisterFunc("delete", memcached.RedisDelete)
	server.RegisterFunc("incr", memcached.RedisIncr)
	server.RegisterFunc("flush_all", memcached.RedisFlushAll)
	server.RegisterFunc("version", memcached.RedisVersion)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
