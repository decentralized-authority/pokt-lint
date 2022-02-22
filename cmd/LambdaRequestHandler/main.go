package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/itsnoproblem/pokt-lint/linting"
	"github.com/pkg/errors"
	"net/http"
)

type LintRequest struct {
	NodeURL string   `json:"node_url"`
	NodeID  string   `json:"node_id"`
	Chains  []string `json:"chain_ids"`
}

type LintResponse struct {
	StatusCode  float64       `json:"status_code"`
	Message     string        `json:"message"`
	PingResult  PingResponse  `json:"ping_result"`
	RelayResult RelayResponse `json:"relay_result"`
}

type PingResponse struct {
	Success    bool    `json:"success"`
	AvgRTT     float64 `json:"avg_rtt_ms"`
	PacketLoss float64 `json:"packet_loss"`
}

type ChainResult struct {
	ChainID    string `json:"chain_id"`
	Duration   string `json:"duration"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type RelayResponse struct {
	Chains map[string]ChainResult `json:"chains"`
}

func LambdaRequestHandler(ctx context.Context, req LintRequest) (LintResponse, error) {
	httpClient := xray.Client(http.DefaultClient)
	if httpClient == nil {
		return LintResponse{}, errors.New("LambdaRequestHandler: failed to initialize http client")
	}

	linter, err := linting.NewNodeChecker(req.NodeID, req.NodeURL, *httpClient)
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	var pingResponse PingResponse
	pingRes, err := linter.RunPingTest(ctx)
	if err != nil {
		pingResponse = PingResponse{
			Success: false,
		}
	} else {
		pingResponse = PingResponse{
			Success:    true,
			AvgRTT:     float64(pingRes.AvgRtt.Milliseconds()),
			PacketLoss: pingRes.PacketLoss,
		}
	}

	relayResult, err := linter.RunRelayTests()
	if err != nil {
		return LintResponse{}, fmt.Errorf("LambdaRequestHandler: %s", err)
	}

	chainResult := make(map[string]ChainResult, len(relayResult))
	for _, ch := range relayResult {
		chainResult[ch.ChainID] = ChainResult{
			ChainID:    ch.ChainID,
			Duration:   "",
			Successful: ch.Successful,
			Message:    ch.Message,
		}
	}

	message := "some message could go here"
	return LintResponse{
		StatusCode:  200,
		Message:     message,
		PingResult:  pingResponse,
		RelayResult: RelayResponse{Chains: chainResult},
	}, nil
}

func main() {
	lambda.Start(LambdaRequestHandler)
}
