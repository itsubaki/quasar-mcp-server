package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/itsubaki/quasar-mcp-server/quasar"
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

func main() {
	s := quasar.NewMCPServer(identityToken, targetURL)
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return s
	}, nil)

	if err := http.ListenAndServe(addr, handler); err != nil {
		panic(err)
	}
}
