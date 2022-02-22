package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/linting"
	"net/http"
	"net/url"
)

const (
	httpClientTimeoutSec = 10
)

type LintRequest struct {
	NodeURL string   `json:"node_url"`
	NodeID  string   `json:"node_id"`
	Chains  []string `json:"chain_ids"`
}

type LintResponse struct {
	StatusCode float64 `json:"status_code"`
	Message    string  `json:"message"`
}

func LambdaRequestHandler(ctx context.Context, req LintRequest) (LintResponse, error) {
	httpClient := http.DefaultClient
	if httpClient == nil {
		return LintResponse{}, errors.New("LambdaRequestHandler: http client was nil")
	}

	linter, err := linting.NewNodeChecker(req.NodeID, req.NodeURL, *httpClient)
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	pingRes, err := linter.RunPingTest(ctx)
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	message, err := json.Marshal(pingRes)
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	return LintResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("ping response: %s", message),
	}, nil
}

func main() {
	ctx := context.Background()
	req := LintRequest{
		NodeURL: "https://node-000.pokt.gaagl.com",
		Chains:  nil,
	}
	response, err := LambdaRequestHandler(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("panic: got error from LambdaRequestHandler: %s", err))
	}

	output, err := json.Marshal(response)
	if err != nil {
		panic(fmt.Sprintf("panic while marshaling respomnse: %s", err))
	}

	fmt.Printf("output: %s", output)
}

func hostnameFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("hostnameFromURL: %s", err)
	}

	return parsed.Host, nil
}
