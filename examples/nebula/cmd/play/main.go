package main

import (
	"algorithms/examples/nebula/internal/assets"
	"algorithms/examples/nebula/internal/config"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

var tpl *template.Template
var cwd string

func init() {
	var err error
	tpl, err = template.ParseFS(assets.FS, "*.tpl")
	if err != nil {
		panic(err)
	}
	// 获取当前工作目录
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

}

func main() {
	peersCfg := config.NewPeers()
	data, err := os.ReadFile("./conf.d/peers.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, peersCfg)
	if err != nil {
		panic(err)
	}

	err = peersCfg.Ensure()
	if err != nil {
		panic(err)
	}

	lighthouses := make([]*config.Peer, 0, len(peersCfg.Peers))
	for _, peer := range peersCfg.Peers {
		if !peer.Lighthouse {
			continue
		}
		lighthouses = append(lighthouses, peer)
	}

	for _, peer := range peersCfg.Peers {
		log.Println(peer)

		// config
		buf := bytes.NewBuffer(nil)
		err = tpl.ExecuteTemplate(buf, "config.yaml.tpl", map[string]any{
			"lighthouses": lighthouses,
			"peer":        peer,
		})
		if err != nil {
			panic(err)
		}
		data := bytes.NewBuffer(nil)
		scanner := bufio.NewScanner(buf)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.TrimSpace(line), "#") {
				continue
			}
			if len(line) == 0 {
				continue
			}
			fmt.Fprintln(data, line)
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		err := writeFile(fmt.Sprintf("./output/%s.yaml", peer.Name), data.Bytes())
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}

		// sign
		err = os.Chdir(path.Join(cwd, "output"))
		if err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("nebula-cert", "sign", "-name", peer.Name, "-ip", peer.IP, "-groups", strings.Join(peer.Groups, ","))
		output, err := cmd.Output()
		if err != nil {
			if ee, ok := err.(*exec.ExitError); ok {
				log.Printf("%s", ee.Stderr)
			}
			panic(err)
		}
		log.Printf("-- %s", output)
		err = os.Chdir(cwd)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func writeFile(filePath string, data []byte) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
