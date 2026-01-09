package resources

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const lexerURL = "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Lexer.g4"

func NewLexer() (*mcp.Resource, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)) {
	return &mcp.Resource{
			Name:        "openqasm3_lexer_grammar",
			Description: "The OpenQASM 3.x Lexer grammar",
			MIMEType:    "text",
			URI:         lexerURL,
		}, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			body, err := HttpGet(lexerURL)
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{
						URI:      lexerURL,
						Text:     string(body),
						MIMEType: "text",
					},
				},
			}, nil
		}
}
