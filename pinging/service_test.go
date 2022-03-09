package pinging_test

import (
	"context"
	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/mock"
	"github.com/itsnoproblem/pokt-lint/pinging"
	"testing"
)

func TestService_PingHost(t *testing.T) {
	ctx, client := setupTests(true)
	url := "https://www.fake.com"
	svc, err := pinging.NewService(client, url)
	if err != nil {
		t.Fatalf("Error instantiating pinging svc: %s", err)
	}
	svc.SetNumPings(ctx, 1)

	if _, err := svc.PingHost(ctx); err != nil {
		t.Fatalf("Got error pinging host: %s", err)
	}
}

func TestService_PingHost_FailsAsExpected(t *testing.T) {
	ctx, client := setupTests(false)
	url := "https://www.fake.com"
	svc, err := pinging.NewService(client, url)
	if err != nil {
		t.Fatalf("Error instantiating pinging svc: %s", err)
	}
	svc.SetNumPings(ctx, 1)

	stats, err := svc.PingHost(ctx)
	if err != nil {
		t.Fatalf("Got error pinging host: %s", err)
	}
	if stats.NumOk > 0 {
		t.Fatalf("Expected ping to fail but num_ok was %d", stats.NumOk)
	}
}

func TestService_SetNumPings(t *testing.T) {
	ctx, client := setupTests(true)
	url := "https://www.fake.com"
	numPings := int64(3)

	svc, err := pinging.NewService(client, url)
	if err != nil {
		t.Fatalf("Error instantiating pinging svc: %s", err)
	}
	svc.SetNumPings(ctx, numPings)

	stats, err := svc.PingHost(ctx)
	if err != nil {
		t.Fatalf("Got error pinging host: %s", err)
	}
	if stats == nil {
		t.Fatalf("Ping response was null")
	}

	if stats.NumSent != numPings {
		t.Fatalf("requested %d pings to be sent, but %d were sent", numPings, stats.NumSent)
	}

	if stats.NumOk != stats.NumSent {
		t.Fatalf("pings sent (%d) did not match pings ok (%d)", stats.NumSent, stats.NumOk)
	}
}

func setupTests(httpReqShouldSucceed bool) (context.Context, http.Client) {
	return context.Background(), mock.NewFakeHTTPClient(httpReqShouldSucceed)
}
