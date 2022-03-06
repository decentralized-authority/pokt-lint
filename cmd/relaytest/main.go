package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/itsnoproblem/pokt-lint/relaying"
)

func main() {
	nodeURL := flag.String("url", "", "node url")
	nodeID := flag.String("id", "", "node id")
	chainsArg := flag.String("chains", "", "comma separated chains ids, eg: -chains=0001,0003,0005")
	flag.Parse()

	if *nodeURL == "" || *nodeID == "" {
		flag.Usage()
		os.Exit(1)
	}

	var chains []string
	if *chainsArg == "" {
		chains = make([]string, 0)
	} else {
		chains = strings.Split(*chainsArg, ",")
	}

	ctx := context.Background()
	req := relaying.RelayTestRequest{
		NodeURL: *nodeURL,
		NodeID:  *nodeID,
		Chains:  chains,
	}
	response, err := relaying.HandleRequest(ctx, req)
	if err != nil {
		fmt.Printf("Error from LambdaRelayTestHandler: %s", err)
		os.Exit(9)
	}

	output, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error while marshaling respomnse: %s", err)
		os.Exit(9)
	}

	fmt.Printf("%s", output)
}
