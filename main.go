package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/itsubaki/quasar/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/api/idtoken"
)

const (
	lexerURL  = "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Lexer.g4"
	parserURL = "https://raw.githubusercontent.com/itsubaki/qasm/refs/heads/main/qasm3Parser.g4"
)

var (
	identityToken = os.Getenv("IDENTITY_TOKEN")
	targetURL     = os.Getenv("TARGET_URL")
	addr          = func() string {
		port := os.Getenv("PORT")
		if port == "" {
			return ":8080"
		}

		return fmt.Sprintf(":%s", port)
	}()
)

func newQuasarClient(ctx context.Context) (*client.Client, error) {
	if identityToken != "" {
		return client.New(targetURL, client.NewWithIdentityToken(identityToken)), nil
	}

	httpClient, err := idtoken.NewClient(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("new quasar client: %w", err)
	}

	return client.New(targetURL, httpClient), nil
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
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

	return body, nil
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
		mcp.NewTool("openqasm3p0_run",
			mcp.WithDescription("run a quantum circuit using OpenQASM 3.0"),
			mcp.WithString("code",
				mcp.Required(),
				mcp.Description("quantum circuit in OpenQASM 3.0 format"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// parameters
			argCode, ok := req.GetArguments()["code"]
			if !ok {
				return nil, fmt.Errorf("missing required argument")
			}

			code, ok := argCode.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type for code(%T) argument", code)
			}

			// client
			client, err := newQuasarClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("new client: %w", err)
			}

			// run quantum circuit
			resp, err := client.Simulate(ctx, code)
			if err != nil {
				return nil, fmt.Errorf("simulate: %w", err)
			}

			// response
			bytes, err := json.MarshalIndent(resp, "", " ")
			if err != nil {
				return nil, fmt.Errorf("marshal indent: %w", err)
			}

			return mcp.NewToolResultText(string(bytes)), nil
		},
	)

	s.AddResource(
		mcp.NewResource(
			lexerURL,
			"openqasm3p0_lexer_grammar",
			mcp.WithResourceDescription("The OpenQASM3.0 Lexer grammar"),
			mcp.WithMIMEType("text"),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			body, err := httpGet(lexerURL)
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      lexerURL,
					Text:     string(body),
					MIMEType: "text",
				},
			}, nil
		},
	)

	s.AddResource(
		mcp.NewResource(
			parserURL,
			"openqasm3p0_parser_grammar",
			mcp.WithResourceDescription("The OpenQASM3.0 Parser grammar"),
			mcp.WithMIMEType("text"),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			body, err := httpGet(parserURL)
			if err != nil {
				return nil, fmt.Errorf("get: %w", err)
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      parserURL,
					Text:     string(body),
					MIMEType: "text",
				},
			}, nil
		},
	)

	// start
	if err := server.NewStreamableHTTPServer(s).Start(addr); err != nil {
		panic(err)
	}
}
