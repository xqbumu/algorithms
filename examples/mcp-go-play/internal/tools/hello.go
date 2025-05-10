package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Add tool handler
func Hello() server.ServerTool {
	return server.ServerTool{
		Tool:    helloTool,
		Handler: helloHandler,
	}
}

// Add helloTool
var helloTool = mcp.NewTool("hello",
	mcp.WithDescription("Say hello to someone"),
	mcp.WithString("name",
		mcp.Required(),
		mcp.Description("Name of the person to greet"),
	),
)

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
