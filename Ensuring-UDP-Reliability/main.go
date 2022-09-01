package main

import (
	"flag"
	"io/ioutil"
	"log"
)

var (
	address = flag.String("a", "127.0.0.1:6060", "listen address")
	payload = flag.String("p", "/home/sxntana/Documents/coding/Go/Network-programming/payload", "file to serve to clients")
)

func main() {
	flag.Parse()

	p, err := ioutil.ReadFile(*payload)
	if err != nil {
		log.Fatal(err)
	}

	s := Server{Payload: p}
	log.Fatal(s.ListenAndServe(*address))
}
