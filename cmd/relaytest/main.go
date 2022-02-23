package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/linting"
)

const (
	testNodeURL     = "https://YOUR-NODE_URL"
	testNodeAddress = "YOUR-NODE-ADDRESS"
)

func main() {
	ctx := context.Background()
	req := linting.LintRequest{
		NodeURL: testNodeURL,
		NodeID:  testNodeAddress,
		Chains:  []string{"0001", "0005", "0021", "0027", "0040"},
	}
	response, err := linting.HandleRequest(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("panic: got error from LambdaRequestHandler: %s", err))
	}

	output, err := json.Marshal(response)
	if err != nil {
		panic(fmt.Sprintf("panic while marshaling respomnse: %s", err))
	}

	fmt.Printf("%s", output)
}
