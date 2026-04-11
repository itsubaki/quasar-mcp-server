package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ValidateInput struct {
	Code string `json:"code" jsonschema:"quantum circuit in OpenQASM format"`
}

type ValidateOutput struct {
	Valid   bool    `json:"valid" jsonschema:"whether the OpenQASM code is valid"`
	Line    *int32  `json:"line,omitempty" jsonschema:"line number of the error, if any"`
	Column  *int32  `json:"column,omitempty" jsonschema:"column number of the error, if any"`
	Message *string `json:"message,omitempty" jsonschema:"error message if the code is invalid"`
}

func NewValidate(identityToken, targetURL string) (
	*mcp.Tool,
	mcp.ToolHandlerFor[*ValidateInput, *ValidateOutput],
) {
	return &mcp.Tool{
			Name:        "openqasm_validate",
			Description: "Validate a quantum circuit using OpenQASM",
		}, func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input *ValidateInput,
		) (
			*mcp.CallToolResult,
			*ValidateOutput,
			error,
		) {
			client, err := NewQuasarClient(ctx, identityToken, targetURL)
			if err != nil {
				return nil, nil, fmt.Errorf("new quasar client: %w", err)
			}

			resp, err := client.Validate(ctx, input.Code)
			if err != nil {
				return nil, nil, fmt.Errorf("validate: %w", err)
			}

			return nil, &ValidateOutput{
				Valid:   resp.Valid,
				Line:    resp.Line,
				Column:  resp.Column,
				Message: resp.Message,
			}, nil
		}
}
