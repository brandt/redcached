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
	server.RegisterFunc("get", redcached.DefaultGet)
	server.RegisterFunc("set", redcached.DefaultSet)
	server.RegisterFunc("delete", redcached.DefaultDelete)
	server.RegisterFunc("incr", redcached.DefaultIncr)
	server.RegisterFunc("flush_all", redcached.DefaultFlushAll)
	server.RegisterFunc("version", redcached.DefaultVersion)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
