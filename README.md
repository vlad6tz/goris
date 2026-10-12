# GORIS

> A tiny TCP server written from scratch to understand what actually happens inside a socket.

`goris` is a **learning project** - not a production-ready datastore (despite the name). The goal is to
build up, layer by layer, a server that speaks a simple key-value protocol, while digging into:

- **Go's concurrency model** - goroutines, channels, `select`, and where the runtime actually
  switches between them
- **Networking below the standard library** - file descriptors, the `accept` loop, blocking reads,
  and why one goroutine per connection is both simple and wrong at scale
- **Protocol design** - framing, parsing, and eventually getting the server to handle real
  `GET`/`SET` requests instead of echoing whatever shows up


## Roadmap

- [x] Accept TCP connections and echo data back
- [x] Per-connection read deadlines so idle clients can't hold a goroutine forever
- [x] Move to a worker pool instead of unbounded goroutine-per-connection
- [/] Define a real line-based protocol (`GET key`, `SET key value`)
- [x] In-memory key-value store with a mutex or sharded map
- [ ] Graceful shutdown on `SIGINT` / `SIGTERM`

## Running it

```bash
go run ./cmd
```

Then, in another terminal:

```bash
nc localhost 8080
hello goris
# Received message
```
