package tools

import (
	"algorithms/examples/mcp/internal/common"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TinyImage() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(string(GET_TINY_IMAGE),
			mcp.WithDescription("Returns the MCP_TINY_IMAGE"),
		),
		Handler: handleGetTinyImageTool,
	}
}

func handleGetTinyImageTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "This is a tiny image:",
			},
			mcp.ImageContent{
				Type:     "image",
				Data:     common.MCP_TINY_IMAGE,
				MIMEType: "image/png",
			},
			mcp.TextContent{
				Type: "text",
				Text: "The image above is the MCP tiny image.",
			},
		},
	}, nil
}
