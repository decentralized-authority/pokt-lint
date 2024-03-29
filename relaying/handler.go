package relaying

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"log"
	nethttp "net/http"
	"time"
)

const (
	httpClientTimeoutSec = 5
)

const (
	maxNumSamples     = 50
	defaultNumSamples = 5
)

// RelayTestRequest represents the request format the relaying service accepts
type RelayTestRequest struct {
	NodeURL    string   `json:"node_url"`
	NodeID     string   `json:"node_id"`
	Chains     []string `json:"chain_ids"`
	NumSamples int64    `json:"num_samples"`
}

// Validate ensures the validity of all required fields
func (req RelayTestRequest) Validate() error {
	if req.NodeID == "" && len(req.Chains) == 0 {
		return fmt.Errorf("you must specify either 'node_id' or 'chain_ids'")
	}

	if req.NumSamples > maxNumSamples {
		return fmt.Errorf("num_samples cannot exceed %d", maxNumSamples)
	}

	return nil
}

// RelayTestResult represents the result of a relay test
type RelayTestResult struct {
	ChainID          string            `json:"chain_id"`
	ChainName        string            `json:"chain_name"`
	Successful       bool              `json:"success"`
	StatusCode       int               `json:"status_code"`
	Message          string            `json:"message"`
	DurationAvgMS    float64           `json:"duration_avg_ms"`
	DurationMedianMS float64           `json:"duration_median_ms"`
	DurationMinMS    float64           `json:"duration_min_ms"`
	DurationMaxMS    float64           `json:"duration_max_ms"`
	RelayRequest     pocket.Payload    `json:"relay_request"`
	RelayResponses   []RelayTestSample `json:"relay_responses"`
}

// RelayTestSample is an atomic relay test result
type RelayTestSample struct {
	DurationMS float64         `json:"duration_ms"`
	StatusCode int             `json:"status_code"`
	Data       json.RawMessage `json:"data"`
}

// RelayTestResponse is a map of chain ID to RelayTestResult
type RelayTestResponse map[string]RelayTestResult

// HandleRequest handles a relaying service request
func HandleRequest(ctx context.Context, req RelayTestRequest) (RelayTestResponse, error) {
	if req.NumSamples == 0 {
		req.NumSamples = defaultNumSamples
	}

	if err := req.Validate(); err != nil {
		return RelayTestResponse{}, fmt.Errorf("request was invalid: %s", err)
	}

	client := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}
	loggingClient := http.NewWebClient(client, log.Default())

	linter, err := NewService(req.NodeID, req.NodeURL, req.Chains, loggingClient)
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	relayRes, err := linter.RunRelayTests(ctx, req.NumSamples)
	if err != nil {
		return RelayTestResponse{}, fmt.Errorf("relaying.HandleRequest: %s", err)
	}

	return relayRes, nil
}
