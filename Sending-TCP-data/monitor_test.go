package main

import (
	"io"
	"log"
	"net"
	"os"
)

// Monitor embeds Logger for logging network traffic
type Monitor struct {
	*log.Logger
}

// Write implements the io.Writer interface
func (m *Monitor) Write(p []byte) (int, error) {
	return len(p), m.Output(2, string(p))
}

func ExampleMonitor() {
	monitor := &Monitor{Logger: log.New(os.Stdout, "monitor: ", 0)}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}

		defer conn.Close()

		b := make([]byte, 1024)

		// This returns a reader (r) that reads from the network connection
		// and writes all input to the monitor before passing the input to
		// the caller
		r := io.TeeReader(conn, monitor)

		n, err := r.Read(b)
		if err != nil && err != io.EOF {
			monitor.Println(err)
			return
		}

		// This creates a writer (w) that duplicates its writes to the provided
		// writers, i.e conn and monitor
		w := io.MultiWriter(conn, monitor)

		_, err = w.Write(b[:n])
		if err != nil && err != io.EOF {
			monitor.Println(err)
			return
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		monitor.Fatal(err)
	}

	_, err = conn.Write([]byte("Test\n"))
	if err != nil {
		monitor.Fatal(err)
	}

	_ = conn.Close()
	<-done

	// Output:
	// monitor: Test
	// monitor: Test
}
