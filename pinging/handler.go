package pinging

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	nethttp "net/http"
	"time"
)

const (
	httpClientTimeoutSec = 20
	maxNumPings          = 30
)

type PingTestRequest struct {
	NodeURL   string `json:"node_url"`
	PingCount int64  `json:"num_pings"`
}

type PingTestResponse *http.PingStats

func HandleRequest(ctx context.Context, req PingTestRequest) (PingTestResponse, error) {
	httpClient := nethttp.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}

	url := fmt.Sprintf("%s/v1", req.NodeURL)
	pingSvc, err := NewService(&httpClient, url)
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
