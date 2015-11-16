package main

import (
	"fmt"
	"github.com/luxuan/go-memcached-server"
	"github.com/luxuan/go-memcached-server/protocol"
)

// This function needs to be registered.
func Test(req *protocol.McRequest, res *protocol.McResponse) error {
	res.Response = "Awesome custom memcached command implement via function!"
	return nil
}

func main() {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf("Panic: %v\n", msg)
		}
	}()

	methods := map[string]memcached.HandlerFn{
		"get": Test,
	}
	srv, err := memcached.NewServer("", methods)
	if err != nil {
		panic(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
