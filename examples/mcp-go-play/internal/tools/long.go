package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Long() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			string(LONG_RUNNING_OPERATION),
			mcp.WithDescription(
				"Demonstrates a long running operation with progress updates",
			),
			mcp.WithNumber("duration",
				mcp.Description("Duration of the operation in seconds"),
				mcp.DefaultNumber(10),
			),
			mcp.WithNumber("steps",
				mcp.Description("Number of steps in the operation"),
				mcp.DefaultNumber(5),
			),
		),
		Handler: handleLongRunningOperationTool,
	}
}

func handleLongRunningOperationTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	progressToken := request.Params.Meta.ProgressToken
	duration, _ := arguments["duration"].(float64)
	steps, _ := arguments["steps"].(float64)
	stepDuration := duration / steps
	server := server.ServerFromContext(ctx)

	for i := 1; i < int(steps)+1; i++ {
		time.Sleep(time.Duration(stepDuration * float64(time.Second)))
		if progressToken != nil {
			server.SendNotificationToClient(
				ctx,
				"notifications/progress",
				map[string]interface{}{
					"progress":      i,
					"total":         int(steps),
					"progressToken": progressToken,
				},
			)
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"Long running operation completed. Duration: %f seconds, Steps: %d.",
					duration,
					int(steps),
				),
			},
		},
	}, nil
}
