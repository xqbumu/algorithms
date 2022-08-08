package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{}
var transport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: defaultTransportDialContext(&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}),
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("please input a https url follow with exe to handshake.")
	}
	req, err := http.NewRequest(http.MethodGet, os.Args[1], nil)
	if err != nil {
		panic(err)
	}
	client.Transport = transport

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("Response body with len: %d", len(body))
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network string, address string) (net.Conn, error) {
		conn, err := dialer.DialContext(ctx, network, address)
		return &Conn{conn}, err
	}
}

// Conn Add Message Record Type check
// Refer: crypto/tls/common.go
type Conn struct {
	net.Conn
}

// Read implements the Conn Read method.
func (c *Conn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	switch b[0] {
	case 22: // recordTypeHandshake
		log.Printf("Read: %d, Len: %d\n", b[0], len(b))
	default:
		log.Printf("Read: %d, Len: %d\n", b[0], len(b))
	}
	return n, err
}

// Write implements the Conn Write method.
func (c *Conn) Write(b []byte) (int, error) {
	switch b[0] {
	case 22: // recordTypeHandshake, now the cipher is nil
		var vers uint16 = uint16(b[1])<<8 | uint16(b[2])
		if VersionTLS10 == vers {
			log.Printf("Vers: %d\n", vers)
			msg := &clientHelloMsg{}
			if ok := msg.unmarshal(b[5:]); !ok { // skip metadata
				panic("can not unmarshal msg")
			}
			// TODO: make the message face
			msg.supportedSignatureAlgorithms = append(msg.supportedSignatureAlgorithms, 0x1010)
			data := msg.marshal()
			m := len(data)
			record := b[:5]
			record[3] = byte(m >> 8)
			record[4] = byte(m)
			b = append(record, data...)
		}
		log.Printf("Write: %d, Len: %d\n", b[0], len(b))
	default:
		log.Printf("Write: %d, Len: %d\n", b[0], len(b))
	}

	return c.Conn.Write(b)
}
