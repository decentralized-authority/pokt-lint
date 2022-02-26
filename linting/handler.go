package linting

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	nethttp "net/http"
	"time"
)

const (
	httpClientTimeoutSec = 20
)

type LintRequest struct {
	NodeURL string   `json:"node_url"`
	NodeID  string   `json:"node_id"`
	Chains  []string `json:"chain_ids"`
}

type LintResponse struct {
	StatusCode  float64                    `json:"status_code"`
	PingResult  *http.PingStats            `json:"ping_result"`
	RelayResult map[string]RelayTestResult `json:"relay_result"`
}

func HandleRequest(ctx context.Context, req LintRequest) (LintResponse, error) {
	httpClient := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}
	linter, err := NewNodeChecker(req.NodeID, req.NodeURL, httpClient)
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	pingRes, err := linter.RunPingTest(ctx)
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	relayRes, err := linter.RunRelayTests()
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	return LintResponse{
		StatusCode:  200,
		PingResult:  pingRes,
		RelayResult: relayRes,
	}, nil
}
