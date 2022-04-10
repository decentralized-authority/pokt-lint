package mock

import (
	nethttp "net/http"

	"github.com/itsnoproblem/pokt-lint/http"
)

type fakeHTTPClient struct {
	returnSuccessResponses bool
}

// NewFakeHTTPClient returns a mock HTTP client
func NewFakeHTTPClient(successResponses bool) http.Client {
	c := fakeHTTPClient{
		returnSuccessResponses: successResponses,
	}
	return &c
}

// ReturnSuccessResponses determines the type of responses the client returns.
// 200 OK if *true*, 500 Internal Server Error if *false*.
func (c *fakeHTTPClient) ReturnSuccessResponses(newValue bool) http.Client {
	c.returnSuccessResponses = newValue
	return c
}

func (c *fakeHTTPClient) Do(req *nethttp.Request) (*nethttp.Response, error) {
	return c.fakeResponse()
}

func (c *fakeHTTPClient) Get(url string) (*nethttp.Response, error) {
	return c.fakeResponse()
}

func (c *fakeHTTPClient) Options(url string) (*nethttp.Response, error) {
	return c.fakeResponse()
}

func (c *fakeHTTPClient) fakeResponse() (*nethttp.Response, error) {
	var statusCode int
	var status string
	if c.returnSuccessResponses {
		statusCode = 200
		status = "OK"
	} else {
		statusCode = 500
		status = "Internal Server Error"
	}

	res := nethttp.Response{
		Status:     status,
		StatusCode: statusCode,
	}
	return &res, nil
}
