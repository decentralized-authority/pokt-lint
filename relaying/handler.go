package relaying

import (
	"context"
	"fmt"
	nethttp "net/http"
	"time"
)

const (
	httpClientTimeoutSec = 20
)

type RelayTestRequest struct {
	NodeURL string   `json:"node_url"`
	NodeID  string   `json:"node_id"`
	Chains  []string `json:"chain_ids"`
}

type RelayTestResult struct {
	ChainID    string                 `json:"chain_id"`
	Successful bool                   `json:"success"`
	Data       map[string]interface{} `json:"data"`
	DurationMS float64                `json:"duration_ms"`
}

type RelayTestResponse map[string]RelayTestResult

func HandleRequest(ctx context.Context, req RelayTestRequest) (RelayTestResponse, error) {
	httpClient := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}
	linter, err := NewNodeChecker(req.NodeID, req.NodeURL, req.Chains, httpClient)
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	relayRes, err := linter.RunRelayTests()
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	return relayRes, nil
}
