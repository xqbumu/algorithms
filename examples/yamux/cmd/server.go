package cmd

// 多路复用
import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

type Server struct {
	Addr string
}

func (s Server) Start() {
	// 建立底层复用连接
	tcpaddr, _ := net.ResolveTCPAddr("tcp4", s.Addr)
	tcpln, _ := net.ListenTCP("tcp4", tcpaddr)

	conn, _ := tcpln.Accept()
	session, _ := yamux.Server(conn, nil)

	id := 0
	for {
		// 建立多个流通路
		stream, err := session.Accept()
		if err == nil {
			id++
			fmt.Println("accept", id)
			go s.Recv(stream, id)
		} else {
			fmt.Println("session over.")
			return
		}
	}
}

func (s Server) Recv(stream net.Conn, id int) {
	for {
		buf := make([]byte, 4)
		n, err := stream.Read(buf)
		if err == nil {
			fmt.Printf("ID: %d, %d, len: %d, %s\n", id, time.Now().UnixMicro(), n, string(buf))
		} else {
			fmt.Printf("ID: %d, %d, over: %v\n", id, time.Now().UnixMicro(), err)
			return
		}
	}
}
