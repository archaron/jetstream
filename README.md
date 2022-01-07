# Jetstream NATS Module for Helium
[![Build Status](https://travis-ci.com/archaron/jetstream.svg?branch=main)](https://travis-ci.com/archaron/jetstream)
[![Report](https://goreportcard.com/badge/github.com/archaron/jetstream)](https://goreportcard.com/report/github.com/archaron/jetstream)
[![GitHub release](https://img.shields.io/github/release/archaron/jetstream.svg)](https://github.com/archaron/jetstream)
![GitHub](https://img.shields.io/github/license/archaron/jetstream.svg?style=popout)

Module provides you with the following things:
- [`*nats.Conn`](https://godoc.org/github.com/nats-io/nats.go#Conn) represents a bare connection to a nats-server. It can send and receive []byte payloads
- [`nats.JetStreamContext`](https://pkg.go.dev/github.com/nats-io/nats.go#JetStreamContext) represents a streaming context. It can Publish and Subscribe to messages within the NATS JetStream cluster.

