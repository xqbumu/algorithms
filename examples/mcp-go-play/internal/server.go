package internal

import (
	"algorithms/examples/mcp-go-play/internal/hooks"
	"algorithms/examples/mcp-go-play/internal/notification"
	"algorithms/examples/mcp-go-play/internal/prompts"
	"algorithms/examples/mcp-go-play/internal/resources"
	"algorithms/examples/mcp-go-play/internal/tools"

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
		tools.QueryWeather(),
		tools.TinyImage(),
	)

	notification.Add(s)

	return s
}
