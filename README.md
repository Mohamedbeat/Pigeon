# Pigeon

Pigeon is a lightweight, experimental pub/sub (publish/subscribe) message broker written in Go. It is inspired by systems like Redis Pub/Sub and Kafka, but aims to provide a simple, easy-to-understand broker for learning, prototyping, or lightweight use cases.

> **Note:** Pigeon is a side project for learning more about Go and the internals of pub/sub servers. It is still under active development. Features and APIs may change. An official client library is planned but not yet available.

## Features

- **TCP-based server:** Listens for client connections and handles pub/sub commands.
- **Pub/Sub model:** Clients can subscribe to topics and publish messages.
- **Round-robin delivery:** Messages published to a topic are delivered to subscribers in a round-robin fashion.
- **Simple protocol:** Supports commands like `CONNECT`, `SUBSCRIBE`, `PUBLISH`, `UNSUBSCRIBE`, and `DISCONNECT`.
- **Thread-safe topic management:** Uses mutexes to safely manage subscribers and messages.

## How It Works

- Clients connect to the server over TCP (default port: `3333`).
- Each client is assigned a unique ID upon connection.
- Clients can subscribe to topics, publish messages, and unsubscribe.
- When a message is published to a topic, it is delivered to one of the subscribers in a round-robin manner.

## Example Protocol

```
CONNECT
SUBSCRIBE <CLIENT_ID> <TOPIC>
PUBLISH <CLIENT_ID> <TOPIC> <MESSAGE>
UNSUBSCRIBE <CLIENT_ID> <TOPIC>
DISCONNECT
```

Server responses include status messages like `CONNECTED`, `SUBSCRIBED`, `PUBLISHED`, `UNSUBSCRIBED`, `DISCONNECTED`, or error messages.

## Getting Started

### Prerequisites
- Go 1.20 or newer

### Running the Server

Clone the repository and run:

```sh
go run main.go
```

The server will start on port `3333` by default.

### Connecting a Client

You can use `telnet` or `nc` (netcat) to interact with the server manually for testing:

```sh
telnet localhost 3333
```

or

```sh
nc localhost 3333
```

Then type commands as described above.

## Roadmap

- [ ] Official client library
- [ ] Message persistence
- [ ] Authentication and authorization
- [ ] Clustering and high availability
- [ ] Monitoring and metrics
- [ ] More protocol commands and features

## Contributing

Contributions are welcome! Please open issues or pull requests to discuss ideas or report bugs.

## License

MIT License

---

**Pigeon** is a work in progress. Stay tuned for updates and new features! 