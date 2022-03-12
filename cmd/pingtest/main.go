package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	gohttp "net/http"
	"os"
	"time"

	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/pinging"
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

	logger := log.Default()
	c := gohttp.Client{
		Timeout: pingTimeoutMS * time.Millisecond,
	}
	client := http.NewClientWithLogger(&c, logger)

	ctx := context.Background()
	svc, err := pinging.NewService(client, *nodeURL)
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
