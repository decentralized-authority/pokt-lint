package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pinging"
	"net/http"
	"time"
)

const pingTimeoutMS = 1500

func main() {
	c := http.Client{Timeout: pingTimeoutMS * time.Millisecond}
	svc, err := pinging.NewService(c, "https://google.com")
	if err != nil {
		panic(err)
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
