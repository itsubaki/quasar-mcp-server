package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	targetURL = os.Getenv("TARGET_URL")
	addr      = func() string {
		port := os.Getenv("PORT")
		if port == "" {
			return ":8080"
		}

		return fmt.Sprintf(":%s", port)
	}()
)

func newQuasarClient(ctx context.Context) (*client.Client, error) {
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

	s.AddTool(
		mcp.NewTool("factorize",
			mcp.WithDescription("factorize a number using shor's algorithm"),
			mcp.WithString("N",
				mcp.Required(),
				mcp.Description("the number to factorize (string representation of an integer)"),
			),
			mcp.WithString("t",
				mcp.Required(),
				mcp.DefaultString("4"),
				mcp.Description("number of precision qubits (default: 4)"),
			),
			mcp.WithString("a",
				mcp.Required(),
				mcp.DefaultString("-1"),
				mcp.Description("coprime number of N (default: -1, which means a random coprime number will be chosen)"),
			),
			mcp.WithString("seed",
				mcp.Required(),
				mcp.DefaultString("0"),
				mcp.Description("PRNG seed (default: 0, which means a random seed will be chosen)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			getParams := func(name ...string) (map[string]int, error) {
				params := make(map[string]int)
				for _, n := range name {
					arg, ok := req.GetArguments()[n]
					if !ok {
						return nil, fmt.Errorf("missing required argument: %v", name)
					}

					str, ok := arg.(string)
					if !ok {
						return nil, fmt.Errorf("invalid type for %v(%T) argument", name, arg)
					}

					v, err := strconv.Atoi(str)
					if err != nil {
						return nil, fmt.Errorf("convert %v to int: %w", name, err)
					}

					params[n] = v
				}

				return params, nil
			}

			// parameters
			params, err := getParams("N", "t", "a", "seed")
			if err != nil {
				return nil, fmt.Errorf("get parameters: %w", err)
			}

			// client
			client, err := newQuasarClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("new client: %w", err)
			}

			N, t, a, seed := params["N"], min(params["t"], 4), params["a"], params["seed"]
			for range 10 {
				// factorization
				resp, err := client.Factorize(ctx, N, t, a, uint64(seed))
				if err != nil {
					return nil, fmt.Errorf("factorize: %w", err)
				}

				if resp.P == 0 || resp.Q == 0 {
					// no factorization found
					continue
				}

				// response
				if resp.Message != nil {
					// somthing went wrong
					return mcp.NewToolResultText(*resp.Message), nil
				}

				// success
				return mcp.NewToolResultText(strings.Join([]string{
					fmt.Sprintf("The prime factorization of %v is %v and %v.", resp.N, resp.P, resp.Q),
					fmt.Sprintf("num of precision qubits=%v, coprime number of N=%v, PRNG seed=%v, measured bitstring=%v, s/r=%v.", resp.T, resp.A, resp.Seed, resp.M, resp.SR),
				}, "")), nil
			}

			// failed
			return mcp.NewToolResultText(strings.Join([]string{
				"Unfortunately, the operation failed.",
				"Because Shor's algorithm is probabilistic, this can occasionally happen.",
				"Please try again.",
			}, " ")), nil
		},
	)

	s.AddResource(
		mcp.NewResource(
			lexerURL,
			"openqasm3p0_lexer",
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
			"openqasm3p0_parser",
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
