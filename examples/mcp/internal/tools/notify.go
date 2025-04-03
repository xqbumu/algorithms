package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Notify() server.ServerTool {
	return server.ServerTool{
		Tool:    mcp.NewTool("notify"),
		Handler: handleSendNotification,
	}
}

func handleSendNotification(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	server := server.ServerFromContext(ctx)

	err := server.SendNotificationToClient(
		ctx,
		"notifications/progress",
		map[string]interface{}{
			"progress":      10,
			"total":         10,
			"progressToken": 0,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "notification sent successfully",
			},
		},
	}, nil
}
