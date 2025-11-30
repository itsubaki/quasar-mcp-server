package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ShareInput struct {
	Code string `json:"code" jsonschema:"quantum circuit in OpenQASM 3.0 format"`
}

type ShareOutput struct {
	ID        string    `json:"id" jsonschema:"unique identifier for the shared circuit"`
	Code      string    `json:"code" jsonschema:"quantum circuit in OpenQASM 3.0 format"`
	CreatedAt time.Time `json:"created_at" jsonschema:"timestamp when the circuit was shared"`
}

func NewShare(identityToken, targetURL string) (
	*mcp.Tool,
	mcp.ToolHandlerFor[*ShareInput, *ShareOutput],
) {
	return &mcp.Tool{
			Name:        "openqasm3p0_share",
			Description: "Share a quantum circuit using OpenQASM 3.0",
		}, func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input *ShareInput,
		) (
			*mcp.CallToolResult,
			*ShareOutput,
			error,
		) {
			client, err := NewQuasarClient(ctx, identityToken, targetURL)
			if err != nil {
				return nil, nil, fmt.Errorf("new quasar client: %w", err)
			}

			resp, err := client.Share(ctx, input.Code)
			if err != nil {
				return nil, nil, fmt.Errorf("share: %w", err)
			}

			return nil, &ShareOutput{
				ID:        resp.ID,
				Code:      resp.Code,
				CreatedAt: resp.CreatedAt,
			}, nil
		}
}
