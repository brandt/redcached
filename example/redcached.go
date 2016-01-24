package main

import (
	"github.com/brandt/redcached"
)

func main() {
	server, err := redcached.NewServer("", nil)
	if err != nil {
		panic(err)
	}

	// register handler
	server.RegisterFunc("get", redcached.RedisGet)
	server.RegisterFunc("add", redcached.RedisSetNX)
	server.RegisterFunc("set", redcached.RedisSet)
	server.RegisterFunc("delete", redcached.RedisDelete)
	server.RegisterFunc("incr", redcached.RedisIncr)
	server.RegisterFunc("flush_all", redcached.RedisFlushAll)
	server.RegisterFunc("version", redcached.RedisVersion)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
