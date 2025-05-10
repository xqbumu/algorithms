package main

import (
	"github.com/mark3labs/mcp-go/client"
)

func main() {
	_, err := client.NewStdioMCPClient("mcp-server", nil)
	if err != nil {
		panic(err)
	}
}
