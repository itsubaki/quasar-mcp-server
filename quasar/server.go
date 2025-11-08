package quasar

import "github.com/modelcontextprotocol/go-sdk/mcp"

func NewMCPServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{
		Name:    "quasar-mcp-server",
		Version: "v0.0.2",
	}, nil)
}
