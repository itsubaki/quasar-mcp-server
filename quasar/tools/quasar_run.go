package tools

import (
	"context"
	"fmt"

	"github.com/itsubaki/quasar/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type OpenQASM3p0RunInput struct {
	Code string `json:"code" jsonschema:"quantum circuit in OpenQASM 3.0 format"`
}

type OpenQASM3p0RunOutput client.States

func NewOpenQASM3p0Run(identityToken, targetURL string) (
	*mcp.Tool,
	mcp.ToolHandlerFor[*OpenQASM3p0RunInput, *OpenQASM3p0RunOutput],
) {
	return &mcp.Tool{
			Name:        "openqasm_run",
			Description: "Run a quantum circuit using OpenQASM 3.x",
		}, func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input *OpenQASM3p0RunInput,
		) (
			*mcp.CallToolResult,
			*OpenQASM3p0RunOutput,
			error,
		) {
			client, err := NewQuasarClient(ctx, identityToken, targetURL)
			if err != nil {
				return nil, nil, fmt.Errorf("new quasar client: %w", err)
			}

			resp, err := client.Simulate(ctx, input.Code)
			if err != nil {
				return nil, nil, fmt.Errorf("simulate: %w", err)
			}

			out := OpenQASM3p0RunOutput(*resp)
			return nil, &out, nil
		}
}
