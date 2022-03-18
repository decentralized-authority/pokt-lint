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
	NodeURL    string   `json:"node_url"`
	NodeID     string   `json:"node_id"`
	Chains     []string `json:"chain_ids"`
	NumSamples int64    `json:"num_samples"`
}

// RelayTestResult represents the result of a relay test
type RelayTestResult struct {
	ChainID        string            `json:"chain_id"`
	ChainName      string            `json:"chain_name"`
	Successful     bool              `json:"success"`
	StatusCode     int               `json:"status_code"`
	Message        string            `json:"message"`
	DurationMS     float64           `json:"duration_ms"`
	RelayRequest   rpc.Payload       `json:"relay_request"`
	RelayResponses []RelayTestSample `json:"relay_responses"`
}

type RelayTestSample struct {
	DurationMS    float64              `json:"duration_ms"`
	RelayResponse pocket.RelayResponse `json:"response"`
}

// RelayTestResponse is a map of chain ID to RelayTestResult
type RelayTestResponse map[string]RelayTestResult

// HandleRequest handles a relaying service request
func HandleRequest(ctx context.Context, req RelayTestRequest) (RelayTestResponse, error) {
	httpClient := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}

	if req.NumSamples == 0 {
		req.NumSamples = 5
	}

	linter, err := NewNodeChecker(req.NodeID, req.NodeURL, req.Chains, &httpClient)
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	relayRes, err := linter.RunRelayTests(ctx, req.NumSamples)
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	return relayRes, nil
}
