package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Notify() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("notify",
			mcp.WithDescription("Sends a progress notification to client"),
			mcp.WithString("message",
				mcp.Description("Notification message content"),
				mcp.Required(),
			),
			mcp.WithNumber("progress",
				mcp.Description("Current progress value"),
				mcp.DefaultNumber(0),
			),
			mcp.WithNumber("total",
				mcp.Description("Total progress steps"),
				mcp.DefaultNumber(10),
			),
		),
		Handler: handleSendNotification,
	}
}

func handleSendNotification(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments

	message, ok1 := arguments["message"].(string)
	progress, ok2 := arguments["progress"].(int)
	total, ok3 := arguments["total"].(int)
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid number arguments")
	}

	server := server.ServerFromContext(ctx)
	log.Println(message)

	err := server.SendNotificationToClient(
		ctx,
		"notifications/progress",
		map[string]interface{}{
			// "message":       message,
			"progress":      progress,
			"total":         total,
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
