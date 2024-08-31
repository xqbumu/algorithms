package main

import (
	"fmt"
	iofs "io/fs"
	"log"
	"net"

	"github.com/hirochachacha/go-smb2"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.1:445")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "USERNAME",
			Password: "PASSWORD",
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	snames, err := s.ListSharenames()
	if err != nil {
		panic(err)
	}

	for _, sname := range snames {
		log.Println(sname)
	}

	fs, err := s.Mount(`opt`)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	matches, err := iofs.Glob(fs.DirFS("."), "*")
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		fmt.Println(match)
	}

	err = iofs.WalkDir(fs.DirFS("."), ".", func(path string, d iofs.DirEntry, err error) error {
		fmt.Println(path, d, err)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
