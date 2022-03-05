package http_test

import (
	"github.com/itsnoproblem/pokt-lint/http"
	gohttp "net/http"
	"testing"
)

type client interface {
	Do(req *gohttp.Request) (*gohttp.Response, error)
	Get(url string) (gohttp.Response, error)
}

type fakeHttpClient struct {
	returnSuccessResponses bool
}

func (c fakeHttpClient) ReturnSuccessResponses(newValue bool) fakeHttpClient {
	c.returnSuccessResponses = newValue
	return c
}

func (c fakeHttpClient) Do(req *gohttp.Request) (*gohttp.Response, error) {
	return c.fakeResponse()
}

func (c fakeHttpClient) Get(url string) (*gohttp.Response, error) {
	return c.fakeResponse()
}

func (c fakeHttpClient) fakeResponse() (*gohttp.Response, error) {
	var statusCode int
	var status string
	if c.returnSuccessResponses {
		statusCode = 200
		status = "OK"
	} else {
		statusCode = 500
		status = "Internal Server Error"
	}

	res := gohttp.Response{
		Status:     status,
		StatusCode: statusCode,
	}
	return &res, nil
}

func TestPinger_RunSuccess(t *testing.T) {
	url := "https://node-000.mynode.com"
	var client http.Client
	client = fakeHttpClient{returnSuccessResponses: true}
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
	var client http.Client
	client = fakeHttpClient{returnSuccessResponses: false}
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
