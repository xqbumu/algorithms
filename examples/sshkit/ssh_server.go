package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
)

func StartSshServer() {
	forwardHandler := &ssh.ForwardedTCPHandler{}
	server := ssh.Server{
		Addr: "127.0.0.1:2222",
		ChannelHandlers: map[string]ssh.ChannelHandler{
			"direct-tcpip": ssh.DirectTCPIPHandler,
			"session":      ssh.DefaultSessionHandler,
		},
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(
			func(ctx ssh.Context, host string, port uint32) bool {
				log.Println("attempt to bind", host, port, "granted")
				return true
			},
		),
		RequestHandlers: map[string]ssh.RequestHandler{
			"tcpip-forward":        forwardHandler.HandleSSHRequest,
			"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
		},
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": sftpHandler,
		},
	}

	server.Handle(forwardAgent)

	log.Fatal(server.ListenAndServe())
}

// sftpHandler handler for SFTP subsystem
func sftpHandler(sess ssh.Session) {
	debugStream := io.Discard
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(debugStream),
	}
	server, err := sftp.NewServer(
		sess,
		serverOptions...,
	)
	if err != nil {
		log.Printf("sftp server init error: %s\n", err)
		return
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
		fmt.Println("sftp client exited session.")
	} else if err != nil {
		fmt.Println("sftp server completed with error:", err)
	}
}

func forwardAgent(s ssh.Session) {
	cmd := exec.Command("ssh-add", "-l")
	if ssh.AgentRequested(s) {
		l, err := ssh.NewAgentListener()
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		go ssh.ForwardAgentConnections(l, s)
		cmd.Env = append(s.Environ(), fmt.Sprintf("%s=%s", "SSH_AUTH_SOCK", l.Addr().String()))
	} else {
		cmd.Env = s.Environ()
	}
	cmd.Stdout = s
	cmd.Stderr = s.Stderr()
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}
}
