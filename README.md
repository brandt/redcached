# redcached

Provides a Memcached protocol interface to Redis for a limited subset of operations.

The proxy server can currently only speak Memcached's ASCII-based protocol.

## Building

From within the repo root:

    make

## Running

    ./redcached

## References

### Source Code

- Protocol and process: https://github.com/zobo/mrproxy
- Redis library: https://godoc.org/gopkg.in/redis.v3

### Memcached Protocol

- https://github.com/memcached/memcached/blob/master/doc/protocol.txt
- http://blog.elijaa.org/2010/05/21/memcached-telnet-command-summary/

## Authors

- J. Brandt Buckley
- Based on: https://github.com/luxuan/go-memcached-server
- Which is itself based upon: https://github.com/docker/go-redis-server
