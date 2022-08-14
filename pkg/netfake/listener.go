package netfake

import "net"

type Listener struct {
	net.Listener
}

func NewListener(ln net.Listener) net.Listener {
	return &Listener{ln}
}

func (ln Listener) Accept() (net.Conn, error) {
	conn, err := ln.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return &Conn{conn}, err
}
