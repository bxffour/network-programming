package echo

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestEchoServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = client.Close() }()

	msg := []byte("PING")
	_, err = client.WriteTo(msg, server)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)

	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != server.String() {
		t.Fatalf("received a reply from %q instead of %q", addr, server)
	}

	if !bytes.Equal(msg, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", msg, buf[:n])
	}
}
