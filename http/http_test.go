package http_test

import (
	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/mock"
	"testing"
)

func TestPinger_RunSuccess(t *testing.T) {
	url := "https://node-000.mynode.com"
	client := mock.NewFakeHTTPClient(true)
	pinger := http.NewPinger(client, url)
	pinger.Count = 3
	stats, err := pinger.Run()
	if err != nil {
		t.Fatalf("got error running pinger: %s", err)
	}

	if stats.NumSent != stats.NumOk {
		t.Fatalf("pings sent (%d) did not match pings ok (%d)", stats.NumSent, stats.NumOk)
	}
}

func TestPinger_RunFailure(t *testing.T) {
	url := "https://node-000.mynode.com"
	client := mock.NewFakeHTTPClient(false)
	pinger := http.NewPinger(client, url)
	pinger.Count = 3
	stats, err := pinger.Run()
	if err != nil {
		t.Fatalf("got error running pinger: %s", err)
	}

	if stats.NumOk > 0 {
		t.Fatalf("expected requests to fail, but %d were ok", stats.NumOk)
	}
}
