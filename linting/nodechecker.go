package linting

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/pinging"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"github.com/itsnoproblem/pokt-lint/timer"
	"net"
	nethttp "net/http"
	"net/url"
)

type NodeChecker interface {
	RunPingTest(ctx context.Context) (*ping.Statistics, error)
	RunRelayTests(ctx context.Context) (map[string]interface{}, error)
}

type nodeChecker struct {
	pinger         pinging.Service
	pocketProvider pocket.Provider
	nodeID         string
	nodeURL        string
	nodeChains     []pocket.Chain
}

type RelayTestResult struct {
	ChainID    string                 `json:"chain_id"`
	Successful bool                   `json:"success"`
	Data       map[string]interface{} `json:"data"`
	DurationMS int64                  `json:"duration_ms"`
}

func NewNodeChecker(nodeID, nodeAddress string, httpClient nethttp.Client) (*nodeChecker, error) {
	empty := nodeChecker{}
	pingSvc, err := pinging.NewService(httpClient, nodeAddress)
	if err != nil {
		return &empty, fmt.Errorf("linting.NewNodeChecker: %s", err)
	}

	pocketProvider := pocket.NewProvider(httpClient, nodeAddress)

	nc := nodeChecker{
		pinger:         pingSvc,
		pocketProvider: pocketProvider,
		nodeID:         nodeID,
		nodeURL:        nodeAddress,
	}

	if err := nc.init(); err != nil {
		return &empty, fmt.Errorf("linting.NewNodeChecker: %s", err)
	}

	return &nc, nil
}

func (c *nodeChecker) RunPingTest(ctx context.Context) (*http.PingStats, error) {
	res, err := c.pinger.PingHost(ctx)
	if err != nil {
		return nil, fmt.Errorf("nodeChecker.RunPingtest: %s", err)
	}

	return res, nil
}

func (c *nodeChecker) RunRelayTests() (map[string]RelayTestResult, error) {
	chains := make(map[string]RelayTestResult, len(c.nodeChains))
	if len(c.nodeChains) < 1 {
		return nil, errors.New(fmt.Sprintf("No chains for node %s", c.nodeID))
	}

	for _, chain := range c.nodeChains {
		var success bool
		msg := make(map[string]interface{})

		t := timer.Start()
		req, err := RPCRequestForChainID(chain.ID)
		if err != nil {
			success = false
			msg = map[string]interface{}{
				"error": fmt.Sprintf("nodeChecker.RunRelayTests: %s", err),
				"code":  500,
			}
		} else {
			res, err := c.pocketProvider.SimulateRelay(req.ChainID, req.Path, req.Payload)
			if err != nil {
				relayErr, ok := err.(pocket.RelayError)
				if !ok {
					relayErr = pocket.NewRelayError(500, err)
				}

				success = false
				msg = map[string]interface{}{
					"error": relayErr.Err,
					"code":  relayErr.Code,
				}
			} else {
				success = true
				msg = res
			}
		}

		chains[req.ChainID] = RelayTestResult{
			ChainID:    req.ChainID,
			Successful: success,
			Data:       msg,
			DurationMS: t.Elapsed().Milliseconds(),
		}
	}
	return chains, nil
}

func (c *nodeChecker) init() error {
	node, err := c.pocketProvider.Servicer(c.nodeID)
	if err != nil {
		return fmt.Errorf("init: %s", err)
	}

	c.nodeURL = node.ServiceURL
	c.nodeChains = node.Chains
	return nil
}

func ipAddressFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("ipAddressFromURL: %s", err)
	}

	addr, err := net.ResolveIPAddr("ip4:1", parsed.Host)
	if err != nil {
		return "", fmt.Errorf("ipAddressFromURL: %s", err)
	}

	return addr.String(), nil
}
