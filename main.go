package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/quasar/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	BaseURL       = os.Getenv("BASE_URL")
	IdentityToken = os.Getenv("IDENTITY_TOKEN")
)

func main() {
	s := server.NewMCPServer(
		"quasar",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddTool(
		mcp.NewTool("factorization",
			mcp.WithDescription("factorize a number using shor's algorithm"),
			mcp.WithNumber("N",
				mcp.Required(),
				mcp.Description("the number to factorize"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// parameters
			N, t, a := int(req.Params.Arguments["N"].(float64)), 3, 2
			var seed uint64

			// factorization
			resp, err := client.
				New(BaseURL, IdentityToken).
				Factorize(ctx, N, t, a, seed)
			if err != nil {
				return nil, fmt.Errorf("factorize: %w", err)
			}

			// response
			msg := fmt.Sprintf("The prime factorization of %v is %v and %v.", resp.N, resp.P, resp.Q)
			if resp.P == 0 || resp.Q == 0 {
				msg = strings.Join([]string{
					"The operation failed.",
					"Please try again.",
					"Since quantum computation rely on probabilistic algorithms, correct results are not always guaranteed.",
				}, "\n")
			}

			return mcp.NewToolResultText(msg), nil
		},
	)

	// start
	if err := server.ServeStdio(s); err != nil {
		panic(err)
	}
}
