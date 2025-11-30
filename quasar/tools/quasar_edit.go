package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EditInput struct {
	ID string `json:"id" jsonschema:"unique identifier for the shared circuit"`
}

type EditOutput struct {
	ID        string    `json:"id" jsonschema:"unique identifier for the shared circuit"`
	Code      string    `json:"code" jsonschema:"quantum circuit in OpenQASM 3.0 format"`
	CreatedAt time.Time `json:"created_at" jsonschema:"timestamp when the circuit was shared"`
}

func NewEdit(identityToken, targetURL string) (
	*mcp.Tool,
	mcp.ToolHandlerFor[*EditInput, *EditOutput],
) {
	return &mcp.Tool{
			Name:        "openqasm3p0_share",
			Description: "Share a quantum circuit using OpenQASM 3.0",
		}, func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input *EditInput,
		) (
			*mcp.CallToolResult,
			*EditOutput,
			error,
		) {
			client, err := NewQuasarClient(ctx, identityToken, targetURL)
			if err != nil {
				return nil, nil, fmt.Errorf("new quasar client: %w", err)
			}

			resp, err := client.Edit(ctx, input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("share: %w", err)
			}

			return nil, &EditOutput{
				ID:        resp.ID,
				Code:      resp.Code,
				CreatedAt: resp.CreatedAt,
			}, nil
		}
}
