package netfake

import (
	"encoding/hex"
	"log"
	"net"
)

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
		log.Printf("Read: %d, Len: %d\n, Content: %s", b[0], len(b), hex.EncodeToString(b))
	default:
		log.Printf("Read: %d, Len: %d, Content: %s\n", b[0], len(b), string(b))
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
