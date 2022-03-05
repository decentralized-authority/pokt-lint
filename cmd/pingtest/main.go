package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pinging"
	"net/http"
	"os"
	"time"
)

const pingTimeoutMS = 1500

func main() {
	numPings := flag.Int64("num", 1, "-num 10")
	nodeURL := flag.String("url", "", "-url https://www.example.com")
	flag.Parse()

	if *nodeURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := http.Client{
		Timeout: pingTimeoutMS * time.Millisecond,
	}

	ctx := context.Background()
	svc, err := pinging.NewService(&c, *nodeURL)
	svc.SetNumPings(ctx, *numPings)
	if err != nil {
		panic(err)
	}

	stats, err := svc.PingHost(ctx)
	if err != nil {
		panic(err)
	}

	encoded, err := json.MarshalIndent(stats, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(encoded))
}
