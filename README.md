## About

This is a bare bones re-implementation of the Memcached server protocol in Go.

## Building

    go build -o server example/default.go
    ./server

## References

### Source Code

- protocol, process: https://github.com/zobo/mrproxy
- main server, client, handler/example: https://github.com/docker/go-redis-server

### Memcached Protocol

- http://acooly.iteye.com/blog/1120346
- http://my.oschina.net/flynewton/blog/10671
