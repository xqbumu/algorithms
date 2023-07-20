package main

import (
	"bytes"
	"io"
	"net"
)

// TeeConn represents a simple net.Conn interface for SNI Processing.
type TeeConn struct {
	Conn      net.Conn
	Reader    io.Reader
	Buffer    *bytes.Buffer
	FirstRead bool
	Flushed   bool
}

func NewTeeConn(conn net.Conn) *TeeConn {
	teeConn := &TeeConn{
		Conn:    conn,
		Buffer:  bytes.NewBuffer([]byte{}),
		Flushed: false,
	}

	teeConn.Reader = io.TeeReader(conn, teeConn.Buffer)

	return teeConn
}

// Read implements a reader ontop of the TeeReader.
func (conn *TeeConn) Read(p []byte) (int, error) {
	if !conn.FirstRead {
		conn.FirstRead = true
		return conn.Reader.Read(p)
	}

	if conn.FirstRead && !conn.Flushed {
		conn.Flushed = true
		copy(p[0:conn.Buffer.Len()], conn.Buffer.Bytes())
		return conn.Buffer.Len(), nil
	}

	return conn.Conn.Read(p)
}

// Write is a shim function to fit net.Conn.
func (conn *TeeConn) Write(p []byte) (int, error) {
	// if !conn.Flushed {
	// 	return 0, io.ErrClosedPipe
	// }

	return conn.Conn.Write(p)
}

// Close is a shim function to fit net.Conn.
func (conn *TeeConn) Close() error {
	if !conn.Flushed {
		return nil
	}

	return conn.Conn.Close()
}
