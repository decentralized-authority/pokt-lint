package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pinging"
)

func main() {
	svc, err := pinging.NewService("google.com")
	if err != nil {
		panic(err);
	}

	stats, err := svc.PingHost(context.Background())
	if err != nil {
		panic(err)
	}

	encoded, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(encoded))
}
