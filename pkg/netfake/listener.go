package netfake

import "net"

type Listener struct {
	net.Listener
	name string
}

func NewListener(name string, ln net.Listener) net.Listener {
	return &Listener{ln, name}
}

func (ln Listener) Accept() (net.Conn, error) {
	conn, err := ln.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return NewConn(ln.name, conn), err
}
