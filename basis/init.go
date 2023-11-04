package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	peerAddr string
	port     int
	help     bool
)

func init() {
	flag.StringVar(&peerAddr, "peer", "", "Peer node address")
	flag.IntVar(&port, "port", 8123, "Node listen port")
	flag.BoolVar(&help, "h", false, "Command help")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: main [--port port] [--peer address]

Options:
`)
	flag.PrintDefaults()
}
