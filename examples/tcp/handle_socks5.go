package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
)

const (
	socks4Version = 0x04
	socks5Version = 0x05

	authNone         = 0x00
	authGaaApi       = 0x01
	authUsernamePass = 0x02
	authChap         = 0x03

	socks4Connect = 0x01
	socks5Connect = 0x01

	socks5AddrTypeIPv4   = 0x01
	socks5AddrTypeDomain = 0x03
	socks5AddrTypeIPv6   = 0x04

	socks4ReplySuccess = 0x5a
	socks4ReplyFail    = 0x5b

	socks5ReplySuccess                = 0x00
	socks5ReplyAuthMethodNotSupported = 0xff
	socks5ReplyAuthFailed             = 0x01
	socks5ReplyCmdNotSupported        = 0x07
	socks5ReplyAddrTypeNotSupported   = 0x08
)

type User struct {
	Username string
	Password string
}

func validateUser(username, password string) (*User, error) {
	// In this example, we only allow the user "test" with password "1234"
	if username == "test" && password == "1234" {
		return &User{Username: username, Password: password}, nil
	}
	return nil, fmt.Errorf("invalid username or password")
}

func handleSocks4Connection(conn net.Conn) {
	defer conn.Close()

	// Read the SOCKS4 request
	buf := make([]byte, 8)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		fmt.Println("Error reading SOCKS4 request:", err)
		return
	}

	if buf[0] != socks4Version || buf[1] != socks4Connect {
		fmt.Println("Unsupported SOCKS4 request:", buf[0], buf[1])
		conn.Write([]byte{0x00, socks4ReplyFail, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	// Parse the destination address
	destAddr := net.IPv4(buf[4], buf[5], buf[6], buf[7]).String()
	destPort := uint16(buf[2])<<8 | uint16(buf[3])

	// Authenticate the user
	_, err = validateUser("", "")
	if err != nil {
		fmt.Println("Error authenticating user:", err)
		conn.Write([]byte{0x00, socks4ReplyFail, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	// Connect to the destination server
	destConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", destAddr, destPort))
	if err != nil {
		fmt.Println("Error connecting to destination server:", err)
		conn.Write([]byte{0x00, socks4ReplyFail, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}
	defer destConn.Close()

	// Send the reply message to the client
	reply := []byte{0x00, socks4ReplySuccess, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err = conn.Write(reply)
	if err != nil {
		fmt.Println("Error sending reply message to client:", err)
		conn.Close()
		return
	}

	// Copy data between the client and the destination server
	go func() {
		_, err = io.Copy(destConn, conn)
		if err != nil {
			fmt.Println("Error copying data from client to destination server:", err)
		}
	}()

	go func() {
		_, err = io.Copy(conn, destConn)
		if err != nil {
			fmt.Println("Error copying data from destination server to client:", err)
		}
	}()
}

func handleSocks5Connection(conn io.ReadWriteCloser) {
	defer conn.Close()

	// Authenticate the user
	err := authenticateUser(conn)
	if err != nil {
		fmt.Println("Error authenticating user:", err)
		return
	}

	// Read the SOCKS5 request
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading SOCKS5 request:", err)
		return
	}
	if buf[0] != socks5Version {
		fmt.Println("Unsupported SOCKS version:", buf[0])
		conn.Write([]byte{socks5Version, socks5ReplyCmdNotSupported, 0x00, socks5AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}
	if buf[1] != socks5Connect {
		fmt.Println("Unsupported SOCKS command:", buf[1])
		conn.Write([]byte{socks5Version, socks5ReplyCmdNotSupported, 0x00, socks5AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	// Parse the destination address
	destAddr, destPort, err := parseDestinationAddress(buf[3:n])
	if err != nil {
		fmt.Println("Error parsing destination address:", err)
		conn.Write([]byte{socks5Version, socks5ReplyAddrTypeNotSupported, 0x00, socks5AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	// Connect to the destination server
	destConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", destAddr, destPort))
	if err != nil {
		fmt.Println("Error connecting to destination server:", err)
		conn.Write([]byte{socks5Version, socks5ReplyAddrTypeNotSupported, 0x00, socks5AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}
	defer destConn.Close()

	// Send the reply message to the client
	reply := []byte{socks5Version, socks5ReplySuccess, 0x00, socks5AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err = conn.Write(reply)
	if err != nil {
		fmt.Println("Error sending reply message to client:", err)
		conn.Close()
		return
	}

	// Copy data between the client and the destination server
	go func() {
		_, err = io.Copy(destConn, conn)
		if err != nil {
			fmt.Println("Error copying data from client to destination server:", err)
		}
	}()

	go func() {
		_, err = io.Copy(conn, destConn)
		if err != nil {
			fmt.Println("Error copying data from destination server to client:", err)
		}
	}()
}

func authenticateUser(conn io.ReadWriter) error {
	// Send the supported authentication methods to the client
	_, err := conn.Write([]byte{socks5Version, 0x02, authNone, authUsernamePass})
	if err != nil {
		return err
	}

	// Read the selected authentication method from the client
	buf := make([]byte, 2)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}
	if buf[0] != socks5Version {
		return fmt.Errorf("unsupported SOCKS version: %d", buf[0])
	}

	switch buf[1] {
	case authNone:
		return nil
	case authUsernamePass:
		// Read the username and password from the client
		buf = make([]byte, 2)
		_, err = conn.Read(buf)
		if err != nil {
			return err
		}
		usernameLen := buf[1]
		usernameBuf := make([]byte, usernameLen)
		_, err = conn.Read(usernameBuf)
		if err != nil {
			return err
		}
		username := string(usernameBuf)

		buf = make([]byte, 2)
		_, err = conn.Read(buf)
		if err != nil {
			return err
		}
		passwordLen := buf[1]
		passwordBuf := make([]byte, passwordLen)
		_, err = conn.Read(passwordBuf)
		if err != nil {
			return err
		}
		password := string(passwordBuf)

		// Validate the username and password
		user, err := validateUser(username, password)
		if err != nil {
			conn.Write([]byte{socks5Version, socks5ReplyAuthFailed})
			return err
		}

		// Authentication successful
		conn.Write([]byte{socks5Version, socks5ReplySuccess})
		fmt.Printf("Authenticated user: %s\n", user.Username)
		return nil
	case authChap:
		// 发送协议版本和支持的身份验证方法
		// 第 1 个字节是 SOCKS 版本号，第 2 个字节是支持的身份验证方法的数量，后面跟着支持的身份验证方法
		_, err = conn.Write([]byte{socks5Version, 2, authNone, authChap})
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// 读取服务器返回的选择的身份验证方法
		// 第 1 个字节是 SOCKS 版本号，第 2 个字节是选择的身份验证方法
		data := make([]byte, 16)
		_, err = conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// 检查服务器选择的身份验证方法是否支持 CHAP
		if bytes.IndexByte(data, authChap) == -1 {
			fmt.Println("CHAP authentication is not supported by the server")
			return nil
		}

		// 生成一个随机的挑战值
		challenge := make([]byte, 16)
		_, err = rand.Read(challenge)
		if err != nil {
			slog.Error("new error", "err", err)
			return err
		}

		// 发送 CHAP 请求
		// 第 1 个字节是 CHAP 版本号，第 2 个字节是挑战值的长度，后面跟着挑战值
		_, err = conn.Write([]byte{1, byte(len(challenge))})
		if err != nil {
			slog.Error("new error", "err", err)
			return err
		}
		_, err = conn.Write(challenge)
		if err != nil {
			slog.Error("new error", "err", err)
			return err
		}

		// 读取服务器的挑战响应
		// 第 1 个字节是 CHAP 版本号，第 2 个字节是响应值的长度，后面跟着响应值
		data = make([]byte, 2+len(challenge))
		_, err = conn.Read(data)
		if err != nil {
			slog.Error("new error", "err", err)
			return err
		}

		// 验证服务器发送的响应值是否正确
		expectedResponse := sha1.Sum(append(challenge, []byte("password")...))
		if !bytes.Equal(data[2:], expectedResponse[:]) {
			slog.Error("new error", "err", fmt.Errorf("incorrect challenge response"))
			return err
		}

		log.Println("xx")
		return nil
	default:
		conn.Write([]byte{socks5Version, socks5ReplyAuthMethodNotSupported})
		return fmt.Errorf("unsupported authentication method: %d", buf[1])
	}
}

func parseDestinationAddress(b []byte) (string, uint16, error) {
	switch b[0] {
	case socks5AddrTypeIPv4:
		if len(b) != 7 {
			return "", 0, fmt.Errorf("invalid IPv4 address length: %d", len(b)-1)
		}
		return net.IPv4(b[1], b[2], b[3], b[4]).String(), uint16(b[5])<<8 | uint16(b[6]), nil
	case socks5AddrTypeDomain:
		domainLen := int(b[1])
		if len(b) != domainLen+3 {
			return "", 0, fmt.Errorf("invalid domain address length: %d", len(b)-3)
		}
		return string(b[2 : domainLen+2]), uint16(b[domainLen+2])<<8 | uint16(b[domainLen+3]), nil
	case socks5AddrTypeIPv6:
		if len(b) != 19 {
			return "", 0, fmt.Errorf("invalid IPv6 address length: %d", len(b)-1)
		}
		ip := net.IP{}
		copy(ip[:], b[1:17])
		return ip.String(), uint16(b[17])<<8 | uint16(b[18]), nil
	default:
		return "", 0, fmt.Errorf("unsupported address type: %d", b[0])
	}
}
