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

	if *nodeURL == "" || *nodeID == "" || *chainsArg == "" {
		flag.Usage()
		os.Exit(1)
	}

	chains := strings.Split(*chainsArg, ",")
	ctx := context.Background()
	req := relaying.RelayTestRequest{
		NodeURL: *nodeURL,
		NodeID:  *nodeID,
		Chains:  chains,
	}
	response, err := relaying.HandleRequest(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("panic: got error from LambdaRelayTestHandler: %s", err))
	}

	output, err := json.Marshal(response)
	if err != nil {
		panic(fmt.Sprintf("panic while marshaling respomnse: %s", err))
	}

	fmt.Printf("%s", output)
}
