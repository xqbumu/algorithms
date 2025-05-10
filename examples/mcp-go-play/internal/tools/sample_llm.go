package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func LLM() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("sampleLLM",
			mcp.WithDescription("Interact with large language models (LLMs) using MCP's sampling feature"),
			mcp.WithString("prompt",
				mcp.Description("The input prompt for the LLM"),
				mcp.Required(),
			),
			mcp.WithNumber("max_tokens",
				mcp.Description("Maximum number of tokens to generate"),
				mcp.DefaultNumber(100),
			),
			mcp.WithNumber("temperature",
				mcp.Description("Controls randomness (0.0-2.0)"),
				mcp.DefaultNumber(1.0),
			),
			mcp.WithNumber("top_p",
				mcp.Description("Nucleus sampling parameter (0.0-1.0)"),
				mcp.DefaultNumber(0.9),
			),
		),
		Handler: handleSampleLLMTool,
	}
}

func handleSampleLLMTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments

	prompt, _ := arguments["prompt"].(string)
	maxTokens, _ := arguments["max_tokens"].(float64)

	// This is a mock implementation. In a real scenario, you would use the server's RequestSampling method.
	result := fmt.Sprintf(
		"Sample LLM result for prompt: '%s' (max tokens: %d)",
		prompt,
		int(maxTokens),
	)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("LLM sampling result: %s", result),
			},
		},
	}, nil
}
