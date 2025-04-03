package notification

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Add(s *server.MCPServer) {
	s.AddNotificationHandler("notification", handleNotification)
}

func handleNotification(
	ctx context.Context,
	notification mcp.JSONRPCNotification,
) {
	log.Printf("Received notification: %s", notification.Method)
}
