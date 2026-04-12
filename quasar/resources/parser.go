package resources

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const ParserURL = "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Parser.g4"

func NewParser() (*mcp.Resource, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)) {
	return &mcp.Resource{
			Name:        "openqasm3_parser_grammar",
			Description: "The OpenQASM 3 Parser grammar",
			MIMEType:    "text",
			URI:         ParserURL,
		}, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			body, err := HttpGet(ParserURL)
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{
						URI:      ParserURL,
						Text:     string(body),
						MIMEType: "text",
					},
				},
			}, nil
		}
}
