package http

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/itsnoproblem/pokt-lint/timer"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

func NewClientWithLogger(c Client, l log.Logger) Client {
	return &ClientWithLogger{
		client: c,
		logger: l,
	}
}

type ClientWithLogger struct {
	client Client
	logger log.Logger
}

func (c ClientWithLogger) logError(msg string, err error) {
	_ = c.logger.Log("type", "ERROR", "loc", log.Caller(2), "error", fmt.Errorf("%s: %s", msg, err))
}

func (c ClientWithLogger) Do(req *http.Request) (*http.Response, error) {
	t := timer.Start()
	resp, err := c.client.Do(req)
	if err != nil {
		c.logError("ClientWithLogger.Do", err)
		return resp, err
	}

	_ = c.logger.Log("type", "INFO", "url", req.URL.String(), "took", t.Elapsed().String())
	return resp, nil
}

func (c ClientWithLogger) Get(url string) (*http.Response, error) {
	t := timer.Start()
	resp, err := c.client.Get(url)
	if err != nil {
		c.logError("ClientWithLogger.Get", err)
		return resp, err
	}
	_ = c.logger.Log("type", "INFO", "url", url, "took", t.Elapsed().String())
	return resp, nil
}

func (c ClientWithLogger) Log(args ...interface{}) error {
	return c.logger.Log(args)
}
