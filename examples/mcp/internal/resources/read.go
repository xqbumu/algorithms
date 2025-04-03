package resources

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddReadResource(s *server.MCPServer) {
	s.AddResource(mcp.NewResource("test://static/resource",
		"Static Resource",
		mcp.WithMIMEType("text/plain"),
	), handleReadResource)
}

func handleReadResource(
	ctx context.Context,
	request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "test://static/resource",
			MIMEType: "text/plain",
			Text:     "This is a sample resource",
		},
	}, nil
}
