package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddCodeReview(s *server.MCPServer) {
	// Code review prompt with embedded resource
	s.AddPrompt(mcp.NewPrompt("code_review",
		mcp.WithPromptDescription("Code review assistance"),
		mcp.WithArgument("pr_number",
			mcp.ArgumentDescription("Pull request number to review"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		prNumber := request.Params.Arguments["pr_number"]
		if prNumber == "" {
			return nil, fmt.Errorf("pr_number is required")
		}

		return mcp.NewGetPromptResult(
			"Code review assistance",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleSystem,
					mcp.NewTextContent("You are a helpful code reviewer. Review the changes and provide constructive feedback."),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewEmbeddedResource(mcp.TextResourceContents{
						URI:      fmt.Sprintf("git://pulls/%s/diff", prNumber),
						MIMEType: "text/x-diff",
					}),
				),
			},
		), nil
	})

}
