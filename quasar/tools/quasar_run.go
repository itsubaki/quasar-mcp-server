package tools

import (
	"context"
	"fmt"

	"github.com/itsubaki/quasar/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RunInput struct {
	Code string `json:"code" jsonschema:"quantum circuit in OpenQASM 3.0 format"`
}

type RunOutput client.States

func NewRun(identityToken, targetURL string) (
	*mcp.Tool,
	mcp.ToolHandlerFor[*RunInput, *RunOutput],
) {
	return &mcp.Tool{
			Name:        "openqasm3_run",
			Description: "Run a quantum circuit using OpenQASM 3.x",
		}, func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input *RunInput,
		) (
			*mcp.CallToolResult,
			*RunOutput,
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

			out := RunOutput(*resp)
			return nil, &out, nil
		}
}
