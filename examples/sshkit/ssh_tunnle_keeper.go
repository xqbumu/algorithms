package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Potterli20/go-flags-fork"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Forward struct {
	SrcHost  string
	SrcPort  int
	DestHost string
	DestPort int
}

type SSHRemote struct {
	Host     string
	Port     int
	Username string
}

var opts struct {
	RemoteForwards []string `short:"R" long:"remote" description:"Forward from remote"`
	LocalForwards  []string `short:"L" long:"local" description:"Forward to remote"`
	Identities     []string `short:"i" long:"identity" description:"Key files"`
	KeepAlive      int      `short:"k" long:"keep-alive" description:"Keep alive interval (seconds)" default:"5"`
	Positional     struct {
		Remote string `description:"Remote host, format: [user]@<host>:[port]" positional-arg-name:"remote" required:"yes"`
	} `positional-args:"yes"`
}

var sshRemote SSHRemote
var localForwards []Forward
var remoteForwards []Forward

func main() {
	args, err := flags.Parse(&opts)

	if err != nil {
		return
	}

	sshRemote, err = parseSSHRemote(opts.Positional.Remote)
	if err != nil {
		log.Fatalf("Failed to parse SSH remote: %v", err)
	}

	for _, str := range opts.LocalForwards {
		forward, err := parseForward(str)
		if err != nil {
			log.Fatalf("Failed to parse local forward: %v", err)
		}
		localForwards = append(localForwards, forward)
	}

	for _, str := range opts.RemoteForwards {
		forward, err := parseForward(str)
		if err != nil {
			log.Fatalf("Failed to parse remote forward: %v", err)
		}
		remoteForwards = append(remoteForwards, forward)
	}

	retryAttempts := 0
	for {
		// Get the SSH agent
		sshAgent, err := getSSHAgent()
		if err != nil {
			log.Printf("Failed to get SSH agent: %v", err)
		}

		config, port, err := getSSHConfig(sshRemote, sshAgent)
		if err != nil {
			log.Printf("Failed to get SSH client config: %v", err)
		}

		if dailSSH(fmt.Sprintf("%s:%d", sshRemote.Host, port), config, args) {
			retryAttempts = 0
		}

		retryDelay := math.Pow(2, float64(retryAttempts)) * 5
		if retryDelay > 300 {
			retryDelay = 300
		}
		log.Printf("Retrying in %d seconds...\n", int(retryDelay))
		time.Sleep(time.Duration(retryDelay) * time.Second)

		retryAttempts++
	}
}

func setupForward(stop chan bool, listener net.Listener, dialer func(src net.Conn) (net.Conn, error)) {
	for {
		srcConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			stop <- true
			return
		}

		destConn, err := dialer(srcConn)
		if err != nil {
			log.Printf("Failed to establish connection: %v", err)
			_ = srcConn.Close()
			continue
		}

		go forwardConnection(srcConn, destConn)
	}
}

func forwardConnection(src net.Conn, dest net.Conn) {
	go pipeTo(src, dest)
	pipeTo(dest, src)
}

func dailSSH(hostPort string, config *ssh.ClientConfig, args []string) bool {
	client, err := ssh.Dial("tcp", hostPort, config)

	if err != nil {
		log.Printf("Failed to connect to SSH server: %v", err)
		return false
	}

	remote := WrapClient(client)

	var listeners []net.Listener

	stop := make(chan bool)

	defer func() {
		_ = remote.Close()
		for _, listener := range listeners {
			_ = listener.Close()
		}
	}()

	for _, opt := range remoteForwards {
		remoteListener, err := remote.ListenAlt(opt.SrcHost, uint32(opt.SrcPort))
		if err != nil {
			log.Printf("Failed to start remote listener: %v", err)
			return true
		}

		listeners = append(listeners, remoteListener)

		log.Printf("Tunnel (R) %s:%d -> (L) %s:%d ", opt.SrcHost, opt.SrcPort, opt.DestHost, opt.DestPort)

		go func(opt Forward) {
			setupForward(stop, remoteListener, func(src net.Conn) (net.Conn, error) {
				localConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opt.DestHost, opt.DestPort))

				if err == nil {
					log.Printf("Froward %s -> (R) %s:%d -> (L) -> %s:%d", src.RemoteAddr(), opt.SrcHost, opt.SrcPort, opt.DestHost, opt.DestPort)
				}

				return localConn, err
			})
		}(opt)
	}

	for _, opt := range localForwards {
		localListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opt.SrcHost, opt.SrcPort))
		if err != nil {
			log.Printf("Failed to start local listener: %v", err)
			return true
		}
		listeners = append(listeners, localListener)

		log.Printf("Tunnel (L) %s:%d -> (R) %s:%d ", opt.SrcHost, opt.SrcPort, opt.DestHost, opt.DestPort)

		go func(opt Forward) {
			setupForward(stop, localListener, func(src net.Conn) (net.Conn, error) {
				remoteConn, err := remote.Dial("tcp", fmt.Sprintf("%s:%d", opt.DestHost, opt.DestPort))

				if err == nil {
					log.Printf("Froward %s -> (L) %s:%d -> (R) -> %s:%d", src.RemoteAddr(), opt.SrcHost, opt.SrcPort, opt.DestHost, opt.DestPort)
				}

				return remoteConn, err
			})
		}(opt)
	}

	if opts.KeepAlive > 0 {
		ticker := time.NewTicker(time.Duration(opts.KeepAlive) * time.Second)
		defer ticker.Stop()

		go func() {
			for {
				<-ticker.C
				_, _, err = remote.SendRequest("keep-alive@golang.org", true, nil)
				if err != nil {
					log.Printf("Failed to send keep-alive: %v", err)
					ticker.Stop()
					stop <- true
					return
				}
			}
		}()
	}

	ses, err := remote.NewSession()

	if err != nil {
		log.Printf("Failed to create session: %v", err)
		return false
	}

	if err == nil {
		stdout, e := ses.StdoutPipe()
		if e != nil {
			log.Printf("Failed setup session IO: %v", e)
			return false
		}
		stderr, e := ses.StderrPipe()
		if e != nil {
			log.Printf("Failed setup session IO: %v", e)
			return false
		}
		if len(args) == 0 {
			err = ses.Shell()
		} else {
			err = ses.Start(strings.Join(args, " "))
		}
		if err != nil {
			log.Printf("Failed to start session: %v", err)
			return false
		}
		go readLine(stdout, "[STDOUT]")
		go readLine(stderr, "[STDERR]")
	} else {
		log.Printf("Failed to start shell for logging: %v", err)
	}

	<-stop

	return true
}

func parseSSHRemote(info string) (SSHRemote, error) {
	// Split the info string into user, host, and port
	parts := strings.SplitN(info, "@", 2)
	res := SSHRemote{Port: -1}
	var hostPort string
	if len(parts) == 1 {
		hostPort = parts[0]
	} else if len(parts) == 2 {
		res.Username = parts[0]
		hostPort = parts[1]
	}

	if strings.Contains(hostPort, ":") {
		host, portStr, err := net.SplitHostPort(hostPort)

		if err != nil {
			return SSHRemote{}, fmt.Errorf("failed to parse host:port: %v", err)
		}

		res.Host = host

		if portStr != "" {
			res.Port, err = strconv.Atoi(portStr)
			if err != nil {
				return SSHRemote{}, fmt.Errorf("failed to parse port: %v", err)
			}
			if res.Port < 1 || res.Port > 65535 {
				return SSHRemote{}, fmt.Errorf("invalid port: %d", res.Port)
			}
		}
	} else {
		res.Host = hostPort
	}

	return res, nil
}

func getSSHAgent() (agent.Agent, error) {
	// Check if SSH_AUTH_SOCK environment variable is set
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock != "" {
		// Connect to SSH agent
		agentConn, err := net.Dial("unix", sshAuthSock)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SSH agent: %v", err)
		}

		// Create SSH agent client
		agentClient := agent.NewClient(agentConn)

		return agentClient, nil
	}

	return nil, nil
}

func getSSHConfig(sshRemote SSHRemote, sshAgent agent.Agent) (*ssh.ClientConfig, int, error) {
	// Check if SSH config file exists
	usr, err := user.Current()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get current user: %v", err)
	}

	var remote = sshRemote

	if remote.Port == -1 {
		portConfig := ssh_config.Get(sshRemote.Host, "Port")
		if portConfig != "" {
			_, err := fmt.Sscanf(portConfig, "%d", &remote.Port)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to parse port: %v", err)
			}
		} else {
			remote.Port = 22
		}
	}

	if remote.Username == "" {
		userConfig := ssh_config.Get(remote.Host, "User")
		if userConfig != "" {
			remote.Username = userConfig
		} else {
			remote.Username = usr.Username
		}
	}

	keyPaths := ssh_config.GetAll(remote.Host, "IdentityFile")

	keyPaths = append(keyPaths, "~/.ssh/id_rsa", "~/.ssh/id_dsa", "~/.ssh/id_ecdsa", "~/.ssh/id_ed25519")

	log.Printf("Connected to: %s@%s:%d", remote.Username, remote.Host, remote.Port)

	cfg, err := createConfig(remote.Username, keyPaths, usr.HomeDir, sshAgent)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create SSH config: %v", err)
	}

	return cfg, remote.Port, nil
}

func createConfig(sshUser string, keyPaths []string, homeDir string, sshAgent agent.Agent) (*ssh.ClientConfig, error) {
	var configSigners []ssh.Signer

	paths := append([]string{}, opts.Identities...)
	paths = append(paths, keyPaths...)

	for _, keyPath := range paths {
		if strings.HasPrefix(keyPath, "~/") {
			keyPath = filepath.Join(homeDir, keyPath[2:])
		}

		// Read the private key file
		privateKey, err := ioutil.ReadFile(keyPath)
		if err != nil {
			log.Printf("failed to read private key file: %v", err)
			continue
		}

		// Create a signer for the private key
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}

		configSigners = append(configSigners, signer)
	}

	if sshAgent != nil {
		signers, err := sshAgent.Signers()
		if err != nil {
			return nil, fmt.Errorf("failed to get configSigners from SSH agent: %v", err)
		}

		signers = append(signers, configSigners...)

		clientConfig := &ssh.ClientConfig{
			User: sshUser,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signers...),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		return clientConfig, nil
	}

	clientConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(configSigners...),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return clientConfig, nil
}

func pipeTo(dst net.Conn, src net.Conn) {
	buf := make([]byte, 0x4000)
	for {
		n, err := src.Read(buf)
		if err != nil {
			break
		}
		_, err = dst.Write(buf[:n])
	}
	_ = src.Close()
	_ = dst.Close()
}

func parseForward(str string) (Forward, error) {
	parts, err := splitParts(str)
	if err != nil {
		return Forward{}, err
	}

	if len(parts) < 2 || len(parts) > 4 {
		return Forward{}, fmt.Errorf("invalid format %s", str)
	}

	fw := Forward{}

	if len(parts) == 2 {
		_, err := fmt.Sscanf(parts[0], "%d", &fw.SrcPort)
		if err != nil {
			return Forward{}, fmt.Errorf("invalid format %s", str)
		}
		_, err = fmt.Sscanf(parts[1], "%d", &fw.DestPort)
		if err != nil {
			return Forward{}, fmt.Errorf("invalid format %s", str)
		}
		fw.SrcHost = "127.0.0.1"
		fw.DestHost = "127.0.0.1"
		return fw, nil
	}

	if len(parts) == 3 {
		_, err := fmt.Sscanf(parts[0], "%d", &fw.SrcPort)
		if err != nil {
			return Forward{}, fmt.Errorf("invalid format %s", str)
		}
		_, err = fmt.Sscanf(parts[2], "%d", &fw.DestPort)
		if err != nil {
			return Forward{}, fmt.Errorf("invalid format %s", str)
		}
		fw.SrcHost = "127.0.0.1"
		fw.DestHost = parts[1]
	}

	_, err = fmt.Sscanf(parts[1], "%d", &fw.SrcPort)
	if err != nil {
		return Forward{}, fmt.Errorf("invalid format %s", str)
	}
	_, err = fmt.Sscanf(parts[3], "%d", &fw.DestPort)
	if err != nil {
		return Forward{}, fmt.Errorf("invalid format %s", str)
	}
	fw.SrcHost = parts[0]
	fw.DestHost = parts[2]

	return fw, nil
}

func splitParts(str string) ([]string, error) {
	var parts []string
	var part string
	var inBrackets bool

	for _, c := range str {
		if c == '[' {
			inBrackets = true
			continue
		} else if c == ']' {
			inBrackets = false
			continue
		} else if c == ':' && !inBrackets {
			parts = append(parts, part)
			part = ""
			continue
		}

		part += string(c)
	}

	if inBrackets {
		return nil, fmt.Errorf("invalid format %s", str)
	}

	if part != "" {
		parts = append(parts, part)
	}

	return parts, nil
}

func readLine(input io.Reader, prefix string) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		log.Printf("%s %s\n", prefix, scanner.Text())
	}
}
