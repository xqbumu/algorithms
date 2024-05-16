package netfake

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

// Conn Add Message Record Type check
// Refer: crypto/tls/common.go
type Conn struct {
	net.Conn
	name  string
	count int
}

func NewConn(name string, conn net.Conn) net.Conn {
	return &Conn{conn, name, 0}
}

// Read implements the Conn Read method.
func (c *Conn) Read(b []byte) (int, error) {
	c.count++
	n, err := c.Conn.Read(b)
	switch b[0] {
	case 22: // recordTypeHandshake
		c.printf("Read: %d, Len: %d\n, Content: %s", b[0], len(b), bytes2hex(b))
	default:
		c.printf("Read: %d, Len: %d, Content: %s\n", b[0], len(b), bytes2hex(b))
	}
	return n, err
}

// Write implements the Conn Write method.
func (c *Conn) Write(b []byte) (int, error) {
	c.count++
	switch b[0] {
	case 22: // recordTypeHandshake, now the cipher is nil
		var vers uint16 = uint16(b[1])<<8 | uint16(b[2])
		if VersionTLS10 == vers {
			c.printf("Vers: %d\n", vers)
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
		c.printf("Write: %d, Len: %d\n", b[0], len(b))
	default:
		c.printf("Write: %d, Len: %d\n", b[0], len(b))
	}

	return c.Conn.Write(b)
}

func (c *Conn) printf(format string, v ...any) {
	prefix := fmt.Sprintf("[%s:%d]", c.name, c.count)
	log.Output(2, fmt.Sprintf(prefix+format, v...))
}

func bytes2hex(b []byte) string {
	bs := append([]byte{}, b...)
	bs = bytes.TrimRight(bs, string([]byte{0}))
	// for bs[len(bs)-1] != '0' {
	// 	bs = bs[0 : len(bs)-1]
	// }
	buf := bytes.NewBuffer([]byte{})
	for k, v := range bs {
		if k%32 == 0 {
			fmt.Fprintf(buf, "\n0x%03X..0x%03X: ", k, k)
		}
		fmt.Fprintf(buf, " 0x%02X", v)
	}

	return buf.String()
}
