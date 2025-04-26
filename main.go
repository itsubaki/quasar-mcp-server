package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/itsubaki/quasar/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	BaseURL       = os.Getenv("BASE_URL")
	IdentityToken = os.Getenv("IDENTITY_TOKEN")
)

func GetIdentityToken() (string, error) {
	if IdentityToken != "" {
		return IdentityToken, nil
	}

	path := os.Getenv("GCLOUD_PATH")
	cmd := exec.Command(path, "auth", "print-identity-token")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("print identity token: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

func main() {
	s := server.NewMCPServer(
		"quasar-mcp-server",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddTool(
		mcp.NewTool("OpenQASMRunner",
			mcp.WithDescription("run a quantum circuit using OpenQASM"),
			mcp.WithString("code",
				mcp.Required(),
				mcp.Description("quantum circuit in OpenQASM format"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// parameters
			code := req.Params.Arguments["code"].(string)

			// auth
			token, err := GetIdentityToken()
			if err != nil {
				return nil, fmt.Errorf("get identity token: %w", err)
			}

			// run quantum circuit
			resp, err := client.
				New(BaseURL, token).
				Run(ctx, code)
			if err != nil {
				return nil, fmt.Errorf("run: %w", err)
			}

			// response
			bytes, err := json.MarshalIndent(resp, "", " ")
			if err != nil {
				return nil, fmt.Errorf("marshal indent: %w", err)
			}

			return mcp.NewToolResultText(string(bytes)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("Factorization",
			mcp.WithDescription("factorize a number using shor's algorithm"),
			mcp.WithNumber("N",
				mcp.Required(),
				mcp.Description("the number to factorize (integer)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// parameters
			N := int(req.Params.Arguments["N"].(float64))
			t, a := 4, -1
			var seed uint64

			// auth
			token, err := GetIdentityToken()
			if err != nil {
				return nil, fmt.Errorf("get identity token: %w", err)
			}

			// factorization
			resp, err := client.
				New(BaseURL, token).
				Factorize(ctx, N, t, a, seed)
			if err != nil {
				return nil, fmt.Errorf("factorize: %w", err)
			}

			// response
			if resp.Message != "" {
				return mcp.NewToolResultText(resp.Message), nil
			}

			if resp.P == 0 || resp.Q == 0 {
				return mcp.NewToolResultText(strings.Join([]string{
					"The operation failed.",
					"Please try again.",
					"Since quantum computation rely on probabilistic algorithms, correct results are not always guaranteed.",
				}, "\n")), nil
			}

			return mcp.NewToolResultText(strings.Join([]string{
				fmt.Sprintf("The prime factorization of %v is %v and %v.", resp.N, resp.P, resp.Q),
				fmt.Sprintf("num of precision qubits=%v, coprime number of N=%v, PRNG seed=%v, measured bitstring=%v, s/r=%v", resp.T, resp.A, resp.Seed, resp.M, resp.SR),
			}, "\n")), nil
		},
	)

	s.AddResource(
		mcp.NewResource(
			"https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Lexer.g4",
			"OpenQASM3Lexer",
			mcp.WithResourceDescription("The OpenQASM3Lexer grammar file"),
			mcp.WithMIMEType("text"),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			resp, err := http.Get("https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Lexer.g4")
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("status code: %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("read all: %w", err)
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Lexer.g4",
					MIMEType: "text",
					Text:     string(body),
				},
			}, nil
		},
	)

	s.AddResource(
		mcp.NewResource(
			"https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Parser.g4",
			"OpenQASM3Parser",
			mcp.WithResourceDescription("The OpenQASM3Parser grammar file"),
			mcp.WithMIMEType("text"),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			resp, err := http.Get("https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Parser.g4")
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("status code: %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("read all: %w", err)
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Parser.g4",
					MIMEType: "text",
					Text:     string(body),
				},
			}, nil
		},
	)

	// start
	if err := server.ServeStdio(s); err != nil {
		panic(err)
	}
}
