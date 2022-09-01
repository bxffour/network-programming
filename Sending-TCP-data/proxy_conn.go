/*
	The aim of this package is to proxy data between a nodes, i.e a source
	and destination node. The proxy acts as an intermediary to connect two
	nodes, i.e the nodes speak to each other through the proxy. The proxy
	relays data sent from one node to the other. In this implementation,
	the io.Copy method used to do this. This method copies infromation from
	the source reader and writes it into the destination writer. It's worth
	noting that the source and destination are both *net.TCPConn objects, the
	data never enters the userspace on linux, thereby causing the data transfer
	to occur more efficiently.
*/
package main

import (
	"io"
	"net"
)

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}

	defer connSource.Close()

	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}

	defer connDestination.Close()

	// connDestination replies to connSource
	go func() { _, _ = io.Copy(connSource, connDestination) }()

	// connSource messages to connDestination
	_, err = io.Copy(connDestination, connSource)

	return err
}
