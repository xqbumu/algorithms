package main

import (
	"algorithms/examples/yamux/cmd"
	"flag"
)

// main is the entry point of the command-line tool.
// It parses the command-line flags and starts the appropriate server or client based on the provided arguments.
// If no valid argument is provided, it prints the default help message.
func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "s", "server":
		cmd.Server{
			Addr: "127.0.0.1:8980",
		}.Start()
	case "c", "client":
		cmd.Client{
			Addr: "127.0.0.1:8980",
		}.Start()
	default:
		println("Usage: yamux [s|server|c|client]")
		println("\ts, server\tStarts the server")
		println("\tc, client\tStarts the client")
	}
}
