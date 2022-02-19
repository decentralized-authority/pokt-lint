package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	urlPathSimulateRelay = "v1/client/sim"
	httpClientTimeoutSec = 10
)

type LintRequest struct {
	NodeURL string   `json:"node_url"`
	Chains  []string `json:"chain_ids"`
}

type LintResponse struct {
	StatusCode int8   `json:"status_code"`
	Message    string `json:"message"`
}

func LambdaRequestHandler(ctx context.Context, name LintRequest) (string, error) {
	response, err := simulateRelay("", "", []byte(``))
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func simulateRelay(servicerUrl, chainID string, payload json.RawMessage) (json.RawMessage, error) {
	url := "https://node-000.pokt.gaagl.com/v1" //fmt.Sprintf("%s/%s", servicerUrl, urlPathSimulateRelay)
	client := http.Client{
		Timeout: httpClientTimeoutSec * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("simulateRelay: %s", err)
	}

	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("doRequest: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doRequest: [%d] %s", resp.StatusCode, err)
	}

	return body, nil
}

func main() {
	lambda.Start(LambdaRequestHandler)
}
