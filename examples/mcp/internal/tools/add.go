package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Add() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(string(ADD),
			mcp.WithDescription("Adds two numbers"),
			mcp.WithNumber("a",
				mcp.Description("First number"),
				mcp.Required(),
			),
			mcp.WithNumber("b",
				mcp.Description("Second number"),
				mcp.Required(),
			),
		),
		Handler: handleAddTool,
	}
}

func handleAddTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	a, ok1 := arguments["a"].(float64)
	b, ok2 := arguments["b"].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid number arguments")
	}
	sum := a + b
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("The sum of %f and %f is %f.", a, b, sum),
			},
		},
	}, nil
}
