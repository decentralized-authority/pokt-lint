package pinging

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	"log"
	nethttp "net/http"
	"time"
)

const (
	httpClientTimeoutSec = 20
	maxNumPings          = 30
)

// PingTestRequest represents the request object for ping tests
type PingTestRequest struct {
	NodeURL   string `json:"node_url"`
	PingCount int64  `json:"num_pings"`
}

// PingTestResponse represent the statistics measured by the ping test
type PingTestResponse *http.PingStats

// HandleRequest handles a pinging service request
func HandleRequest(ctx context.Context, req PingTestRequest) (PingTestResponse, error) {
	httpClient := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}
	client := http.NewClientWithLogger(&httpClient, log.Default())
	url := fmt.Sprintf("%s/v1", req.NodeURL)
	pingSvc, err := NewService(client, url)
	if err != nil {
		return nil, fmt.Errorf("pinging.HandleRequest: %s", err)
	}

	if req.PingCount > maxNumPings {
		return nil, fmt.Errorf("num_pings cannot be greater than %d", maxNumPings)
	}

	if req.PingCount > 0 {
		pingSvc.SetNumPings(ctx, req.PingCount)
	}

	stats, err := pingSvc.PingHost(ctx)
	return stats, nil
}
