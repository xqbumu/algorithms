package config

import (
	"errors"
	"fmt"
	"net"
)

type Cert struct {
	Name string
}

type Peer struct {
	Name       string
	IP         string
	Groups     []string
	Lighthouse bool
	Addr       string `yaml:"addr"`
	Port       int    `yaml:"port"`
}

func (p Peer) IPAddr() string {
	parsedIP, _, err := net.ParseCIDR(p.IP)
	if err != nil {
		panic(err)
	}
	return parsedIP.String()
}

func (p *Peer) Ensure() error {
	var err error

	if !isValidCIDR(p.IP) {
		err = errors.Join(err, fmt.Errorf("ip(%v) is illegal", p.IP))
	}

	if p.Lighthouse {
		if len(p.Addr) == 0 {
			err = errors.Join(err, fmt.Errorf("Addr(%v) is empty", p.Addr))
		}
		if !isValidIP(p.Addr) {
			err = errors.Join(err, fmt.Errorf("Addr(%v) is invalid", p.Addr))
		}
		if p.Port > 60036 || p.Port < 0 {
			err = errors.Join(err, fmt.Errorf("Addr(%v) is illegal", p.Port))
		}
		if p.Port == 0 {
			p.Port = 6262
		}
	}

	return err
}

type Peers struct {
	Cert  Cert
	Peers []*Peer
}

func NewPeers() *Peers {
	return &Peers{
		Peers: make([]*Peer, 0, 10),
	}
}

func (ps *Peers) Ensure() error {
	var err error
	for _, peer := range ps.Peers {
		err = errors.Join(err, peer.Ensure())
	}
	return err
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func isValidCIDR(ip string) bool {
	parsedIP, ipnet, err := net.ParseCIDR(ip)
	return parsedIP != nil && ipnet != nil && err == nil
}
