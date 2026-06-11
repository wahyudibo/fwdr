package main

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestForward(t *testing.T) {
	// net.Pipe gives two in-memory connected conns
	// forward(client, remote): io.Copy(client, remote) reads from remote, writes to client
	client, clientOther := net.Pipe()
	remote, remoteOther := net.Pipe()

	go forward(client, remote)

	want := []byte("hello from remote")
	go func() {
		_, _ = remoteOther.Write(want)
		_ = remoteOther.Close()
	}()

	got := make([]byte, len(want))
	if _, err := io.ReadFull(clientOther, got); err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("got %q, want %q", got, want)
	}
	_ = clientOther.Close()
}

func TestHandleConnection_proxiesData(t *testing.T) {
	// Spin up a local TCP server acting as the remote
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close() //nolint:errcheck

	want := []byte("hello from remote")
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		_, _ = conn.Write(want)
		_ = conn.Close()
	}()

	client, clientOther := net.Pipe()
	dialTimeout = time.Second
	go handleConnection(client, ln.Addr().String(), "test")

	got := make([]byte, len(want))
	if _, err := io.ReadFull(clientOther, got); err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("got %q, want %q", got, want)
	}
	_ = clientOther.Close()
}

func TestHandleConnection_dialFailure(t *testing.T) {
	client, clientOther := net.Pipe()
	defer clientOther.Close() //nolint:errcheck

	dialTimeout = 100 * time.Millisecond
	// Nothing listens on port 1 — dial must fail
	handleConnection(client, "127.0.0.1:1", "test")

	// client must be closed after failed dial
	_, err := clientOther.Write([]byte("x"))
	if err == nil {
		t.Error("expected write to closed connection to fail")
	}
}
