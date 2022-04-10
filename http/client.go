package http

import (
	"fmt"
	"log"
	nethttp "net/http"
	"time"

	"github.com/itsnoproblem/pokt-lint/timer"
)

// Client is an HTTP client.
type Client interface {
	Do(req *nethttp.Request) (*nethttp.Response, error)
	Get(url string) (*nethttp.Response, error)
	Options(url string) (*nethttp.Response, error)
}

// NewWebClient returns a client that writes log messages
func NewWebClient(c nethttp.Client, l *log.Logger) Client {
	return &webClient{
		client: c,
		logger: l,
	}
}

type webClient struct {
	client nethttp.Client
	logger *log.Logger
}

func (c *webClient) Do(req *nethttp.Request) (*nethttp.Response, error) {
	t := timer.Start()
	resp, err := c.client.Do(req)
	if err != nil {
		c.logError("webClient.Do", err)
		return resp, err
	}

	c.logInfo(fmt.Sprintf("%s - took %s", req.URL.String(), t.Elapsed().String()))
	return resp, nil
}

func (c *webClient) Get(url string) (*nethttp.Response, error) {
	t := timer.Start()
	resp, err := c.client.Get(url)
	if err != nil {
		c.logError("webClient.Get", err)
		return resp, err
	}
	c.logInfo(fmt.Sprintf("%s: %s", url, t.Elapsed().String()))
	return resp, nil
}

func (c *webClient) Options(url string) (*nethttp.Response, error) {
	req, err := nethttp.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *webClient) logError(msg string, err error) {
	c.logger.Printf("ERROR - %s - %s: %s", time.Now().String(), msg, err.Error())
}

func (c *webClient) logInfo(msg string) {
	c.logger.Printf("INFO - %s - %s", time.Now().String(), msg)
}
