package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/itsubaki/quasar-mcp-server/quasar"
	"github.com/itsubaki/quasar-mcp-server/quasar/resources"
	"github.com/itsubaki/quasar-mcp-server/quasar/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	identityToken = os.Getenv("IDENTITY_TOKEN")
	targetURL     = os.Getenv("TARGET_URL")
	addr          = func() string {
		port := os.Getenv("PORT")
		if port == "" {
			return ":8080"
		}

		return fmt.Sprintf(":%s", port)
	}()
)

func NewMCPServer(identityToken, targetURL string) *mcp.Server {
	s := quasar.NewMCPServer()

	// resources
	s.AddResource(resources.NewLexer())
	s.AddResource(resources.NewParser())

	// tools
	tool, handler := tools.NewOpenQASM3p0Run(identityToken, targetURL)
	mcp.AddTool(s, tool, handler)

	return s
}

func main() {
	if err := http.ListenAndServe(addr, mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return NewMCPServer(identityToken, targetURL)
	}, nil)); err != nil {
		panic(err)
	}
}
