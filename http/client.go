package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/itsnoproblem/pokt-lint/timer"
)

// Client is an HTTP client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

// NewClientWithLogger returns a client that writes log messages
func NewClientWithLogger(c Client, l *log.Logger) Client {
	return &clientWithLogger{
		client: c,
		logger: *l,
	}
}

type clientWithLogger struct {
	client Client
	logger log.Logger
}

func (c *clientWithLogger) Do(req *http.Request) (*http.Response, error) {
	t := timer.Start()
	resp, err := c.client.Do(req)
	if err != nil {
		c.logError("clientWithLogger.Do", err)
		return resp, err
	}

	c.logInfo(fmt.Sprintf("%s - took %s", req.URL.String(), t.Elapsed().String()))
	return resp, nil
}

func (c *clientWithLogger) Get(url string) (*http.Response, error) {
	t := timer.Start()
	resp, err := c.client.Get(url)
	if err != nil {
		c.logError("clientWithLogger.Get", err)
		return resp, err
	}
	c.logInfo(fmt.Sprintf("%s: %s", url, t.Elapsed().String()))
	return resp, nil
}

func (c *clientWithLogger) logError(msg string, err error) {
	c.logger.Printf("ERROR - %s - %s: %s", time.Now().String(), msg, err.Error())
}

func (c *clientWithLogger) logInfo(msg string) {
	c.logger.Printf("INFO - %s - %s", time.Now().String(), msg)
}
