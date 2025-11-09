package quasar

import (
	"github.com/itsubaki/quasar-mcp-server/quasar/resources"
	"github.com/itsubaki/quasar-mcp-server/quasar/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewMCPServer(identityToken, targetURL string) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "quasar-mcp-server",
		Version: "v0.0.2",
	}, nil)

	// resources
	s.AddResource(resources.NewLexer())
	s.AddResource(resources.NewParser())

	// tools
	tool, handler := tools.NewOpenQASM3p0Run(identityToken, targetURL)
	mcp.AddTool(s, tool, handler)

	return s
}
