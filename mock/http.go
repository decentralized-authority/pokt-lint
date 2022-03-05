package mock

import (
	gohttp "net/http"
)

type fakeHttpClient struct {
	returnSuccessResponses bool
}

func NewFakeHTTPClient(successResponses bool) fakeHttpClient {
	return fakeHttpClient{returnSuccessResponses: successResponses}
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
