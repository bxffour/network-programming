/*
	This package aims to demonstrate why you need to verify the
	sender's address for udp transmissions. In the file, we start
	an echo server, a client and an interloper. The server echoes
	whatever the client sends back, before the client receives a
	reply from the server, the interloper writes to the client. When
	the client reads from its receive buffer it reads the interloper's
	message first then the echo server. In UDP, multiple nodes can send
	data, that's why you need to verify if the sending node's address.
*/
package echo

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestListenPacketUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer cancel()

	// starting the client
	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = client.Close() }()

	// starting the interloper
	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// interloper writing to the client
	interupt := []byte("pardon me")
	n, err := interloper.WriteTo(interupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_ = interloper.Close()

	if l := len(interupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	// client sending a ping message to the server. The server
	// echoes this message back to the client.
	ping := []byte("ping")
	_, err = client.WriteTo(ping, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	// the client reading from its buffer. This would be the interloper's
	// message since it wrote to the client first.
	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(interupt, buf[:n]) {
		t.Errorf("expected reply %q, actual reply %q", interupt, buf[:n])
	}

	// checking if the messages was indeed received from the interloper
	if addr.String() != interloper.LocalAddr().String() {
		t.Errorf("expected message from %q actual sender is %q", interloper.LocalAddr(), addr)
	}

	// reading the server's message from the buffer.
	n, addr, err = client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected reply %q, actual reply %q", ping, buf[:n])
	}

	// verifying if the message was really sent from the server's address
	if addr.String() != serverAddr.String() {
		t.Errorf("expected message from %q, actual sender is %q", serverAddr, addr)
	}
}
