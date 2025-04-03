package internal

import (
	"algorithms/examples/mcp/internal/hooks"
	"algorithms/examples/mcp/internal/notification"
	"algorithms/examples/mcp/internal/prompts"
	"algorithms/examples/mcp/internal/resources"
	"algorithms/examples/mcp/internal/tools"

	"github.com/mark3labs/mcp-go/server"
)

func NewMCPServer() *server.MCPServer {
	s := server.NewMCPServer(
		"example-servers/everything",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks.GetBasic()),
	)

	resources.AddReadResource(s)
	resources.AddTemplate(s)

	prompts.AddSimple(s)
	prompts.AddComplex(s)

	s.AddTools(
		tools.Add(),
		tools.Calculator(),
		tools.Echo(),
		tools.Hello(),
		tools.HTTPRequest(),
		tools.LLM(),
		tools.Long(),
		tools.Notify(),
		tools.TinyImage(),
	)

	notification.Add(s)

	return s
}
