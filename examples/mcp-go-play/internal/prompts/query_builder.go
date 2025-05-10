package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddQueryBuilder(s *server.MCPServer) {
	// Database query builder prompt
	s.AddPrompt(mcp.NewPrompt("query_builder",
		mcp.WithPromptDescription("SQL query builder assistance"),
		mcp.WithArgument("table",
			mcp.ArgumentDescription("Name of the table to query"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		tableName := request.Params.Arguments["table"]
		if tableName == "" {
			return nil, fmt.Errorf("table name is required")
		}

		return mcp.NewGetPromptResult(
			"SQL query builder assistance",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleSystem,
					mcp.NewTextContent("You are a SQL expert. Help construct efficient and safe queries."),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewEmbeddedResource(mcp.TextResourceContents{
						URI:      fmt.Sprintf("db://schema/%s", tableName),
						MIMEType: "application/json",
					}),
				),
			},
		), nil
	})
}
