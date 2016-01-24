package main

import (
	"github.com/brandt/redcached/rcdaemon"
)

func main() {
	server, err := rcdaemon.NewServer("", nil)
	if err != nil {
		panic(err)
	}

	// register handler
	server.RegisterFunc("get", rcdaemon.RedisGet)
	server.RegisterFunc("add", rcdaemon.RedisSetNX)
	server.RegisterFunc("set", rcdaemon.RedisSet)
	server.RegisterFunc("delete", rcdaemon.RedisDelete)
	server.RegisterFunc("incr", rcdaemon.RedisIncr)
	server.RegisterFunc("flush_all", rcdaemon.RedisFlushAll)
	server.RegisterFunc("version", rcdaemon.RedisVersion)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
