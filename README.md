## About

Provides a Memcached protocol interface to Redis for a limited subset of operations.

The proxy server can currently only speak Memcached's ASCII-based protocol.

## Building

    go build -o server example/redis.go
    ./server

## References

### Source Code

- protocol, process: https://github.com/zobo/mrproxy
- main server, client, handler/example: https://github.com/docker/go-redis-server
- https://godoc.org/gopkg.in/redis.v3

### Memcached Protocol

- https://github.com/memcached/memcached/blob/master/doc/protocol.txt
- http://acooly.iteye.com/blog/1120346
- http://my.oschina.net/flynewton/blog/10671
