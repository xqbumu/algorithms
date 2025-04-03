package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func LLM() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.Tool{
			Name:        string(SAMPLE_LLM),
			Description: "Samples from an LLM using MCP's sampling feature",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"prompt": map[string]interface{}{
						"type":        "string",
						"description": "The prompt to send to the LLM",
					},
					"maxTokens": map[string]interface{}{
						"type":        "number",
						"description": "Maximum number of tokens to generate",
						"default":     100,
					},
				},
			},
		},
		Handler: handleSampleLLMTool,
	}
}

func handleSampleLLMTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments

	prompt, _ := arguments["prompt"].(string)
	maxTokens, _ := arguments["maxTokens"].(float64)

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
