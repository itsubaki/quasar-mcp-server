package resources

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const parserURL = "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Parser.g4"

func NewParser() (*mcp.Resource, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)) {
	return &mcp.Resource{
			Name:        "openqasm3_parser_grammar",
			Description: "The OpenQASM 3.x Parser grammar",
			MIMEType:    "text",
			URI:         parserURL,
		}, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			body, err := HttpGet(parserURL)
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{
						URI:      parserURL,
						Text:     string(body),
						MIMEType: "text",
					},
				},
			}, nil
		}
}
