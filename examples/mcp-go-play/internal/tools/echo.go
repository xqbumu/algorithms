package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Add tool handler
func Echo() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(string(ECHO),
			mcp.WithDescription("Echoes back the input"),
			mcp.WithString("message",
				mcp.Description("Message to echo"),
				mcp.Required(),
			),
		),
		Handler: handleEchoTool,
	}
}

func handleEchoTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	message, ok := arguments["message"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid message argument")
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Echo: %s", message),
			},
		},
	}, nil
}
