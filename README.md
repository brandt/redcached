# redcached

Provides a Memcached protocol interface to Redis for a limited subset of operations.

The proxy server can currently only speak Memcached's ASCII-based protocol.

**NOTE: This is a prototype.**

## Building

From within the repo root:

    make

## Running

    ./redcached

## Completeness

Support is mostly complete for the following operations:

- `SET`
- `GET`
- `GETS`
- `ADD`
- `INCR`
- `DECR`
- `FLUSH_ALL`
- `DELETE`

## References

### Source Code

- Protocol and process: https://github.com/zobo/mrproxy
- Redis library: https://godoc.org/gopkg.in/redis.v3

### Memcached Protocol

- https://github.com/memcached/memcached/blob/master/doc/protocol.txt
- http://blog.elijaa.org/2010/05/21/memcached-telnet-command-summary/
- https://github.com/facebook/mcrouter/blob/4d5f15c2f1d2a83c9f0befa30df0923246c9aedb/mcrouter/lib/network/test/MockMc.h

## Authors

- J. Brandt Buckley
- Based on: https://github.com/luxuan/go-memcached-server
- Which is itself based upon: https://github.com/docker/go-redis-server
