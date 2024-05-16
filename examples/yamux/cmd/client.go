package cmd

import (
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

type Client struct {
	Addr string
}

func (c Client) Start() {
	// 建立底层复用通道
	conn, _ := net.Dial("tcp4", c.Addr)
	defer conn.Close()
	session, _ := yamux.Client(conn, nil)
	defer session.Close()

	// 建立应用流通道1
	stream1, _ := session.Open()
	defer stream1.Close()
	stream1.Write([]byte("ping"))
	stream1.Write([]byte("pong"))
	time.Sleep(1 * time.Second)

	// 建立应用流通道2
	stream2, _ := session.Open()
	defer stream2.Close()
	stream2.Write([]byte("pong"))
	time.Sleep(1 * time.Second)

	// 清理退出
	time.Sleep(5 * time.Second)
}
