package relaying

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"github.com/itsnoproblem/pokt-lint/rpc"
	nethttp "net/http"
	"time"
)

const (
	httpClientTimeoutSec = 20
)

// RelayTestRequest represents the request format the relaying service accepts
type RelayTestRequest struct {
	NodeURL string   `json:"node_url"`
	NodeID  string   `json:"node_id"`
	Chains  []string `json:"chain_ids"`
}

// RelayTestResult represents the result of a relay test
type RelayTestResult struct {
	ChainID       string               `json:"chain_id"`
	Successful    bool                 `json:"success"`
	StatusCode    int                  `json:"status_code"`
	Message       string               `json:"message"`
	DurationMS    float64              `json:"duration_ms"`
	RelayRequest  rpc.Payload          `json:"relay_request"`
	RelayResponse pocket.RelayResponse `json:"relay_response"`
}

// RelayTestResponse is a map of chain ID to RelayTestResult
type RelayTestResponse map[string]RelayTestResult

// HandleRequest handles a relaying service request
func HandleRequest(ctx context.Context, req RelayTestRequest) (RelayTestResponse, error) {
	httpClient := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}
	linter, err := NewNodeChecker(req.NodeID, req.NodeURL, req.Chains, &httpClient)
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	relayRes, err := linter.RunRelayTests()
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	return relayRes, nil
}
