package main

import (
	"fmt"
	"github.com/brandt/redcached"
	"github.com/brandt/redcached/protocol"
)

// This function needs to be registered.
func Test(req *protocol.McRequest, res *protocol.McResponse) error {
	res.Response = "Awesome custom redcached command implement via function!"
	return nil
}

func main() {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf("Panic: %v\n", msg)
		}
	}()

	methods := map[string]redcached.HandlerFn{
		"get": Test,
	}
	srv, err := redcached.NewServer("", methods)
	if err != nil {
		panic(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
