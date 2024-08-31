package main

import (
	"context"
	"net"
)

func main() {
	dailer := &net.Dialer{}
	conn, err := dailer.DialContext(context.TODO(), "tcp", "[::ffff:198.19.249.169]:443")
	if err != nil {
		panic(err)
	}

	if err = conn.Close(); err != nil {
		panic(err)
	}
}
