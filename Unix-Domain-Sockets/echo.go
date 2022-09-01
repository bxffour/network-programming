package echo

import (
	"context"
	"net"
	"os"
)

func steamingEchoServer(ctx context.Context, network, address string) (net.Addr, error) {
	// spinning up listener
	s, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		// spinning up infinite loop for accepting connections
		for {
			conn, err := s.Accept()
			if err != nil {
				return
			}

			// a goroutine is spun up to handle each individual connection
			go func() {
				defer func() { _ = conn.Close() }()

				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	// this returns immediately since the code preceeding it runs in a goroutine.
	return s.Addr(), nil
}

func datagramEchoServer(ctx context.Context, network, address string) (net.Addr, error) {
	s, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()

			if network == "unixgram" {
				os.Remove(address)
			}
		}()

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
