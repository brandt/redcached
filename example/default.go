package main

import (
	"github.com/luxuan/go-memcached-server"
)

func main() {
	server, err := memcached.NewServer("", nil)
	if err != nil {
		panic(err)
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
