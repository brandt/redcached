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
	server.RegisterFunc("get", rcdaemon.GetHandler)
	server.RegisterFunc("gets", rcdaemon.GetHandler)
	server.RegisterFunc("add", rcdaemon.AddHandler)
	server.RegisterFunc("set", rcdaemon.SetHandler)
	server.RegisterFunc("delete", rcdaemon.DeleteHandler)
	server.RegisterFunc("incr", rcdaemon.IncrHandler)
	server.RegisterFunc("flush_all", rcdaemon.FlushAllHandler)
	server.RegisterFunc("version", rcdaemon.VersionHandler)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
