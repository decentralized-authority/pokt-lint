package linting

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/itsnoproblem/pokt-lint/pinging"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"net/http"
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
	ChainID    string
	Successful bool
	Message    string
}

func NewNodeChecker(nodeID, nodeAddress string, httpClient http.Client) (nodeChecker, error) {
	host, err := hostnameFromURL(nodeAddress)
	if err != nil {
		return nodeChecker{}, fmt.Errorf("linting.NewNodeChecker: %s", err)
	}

	pingSvc, err := pinging.NewService(host)
	if err != nil {
		return nodeChecker{}, fmt.Errorf("linting.NewNodeChecker: %s", err)
	}

	pocketProvider := pocket.NewProvider(httpClient, nodeAddress)

	nc := nodeChecker{
		pinger:         pingSvc,
		pocketProvider: pocketProvider,
		nodeID:         nodeID,
		nodeURL:        nodeAddress,
	}

	if err := nc.init(); err != nil {
		return nodeChecker{}, fmt.Errorf("linting.NewNodeChecker: %s", err)
	}

	return nc, nil
}

func (c *nodeChecker) RunPingTest(ctx context.Context) (*ping.Statistics, error) {
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
		req, err := RPCRequestForChainID(chain.ID)
		if err != nil {
			return nil, fmt.Errorf("nodeChecker.RunRelayTests: %s", err)
		}

		res, err := c.pocketProvider.SimulateRelay(req.ChainID, req.Path, req.Payload)
		if err != nil {
			return nil, fmt.Errorf("nodeChecker.RunRelayTests: %s", err)
		}

		chains[req.ChainID] = RelayTestResult{
			ChainID:    req.ChainID,
			Successful: true,
			Message:    string(res),
		}
	}
	return nil, nil
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

func hostnameFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("hostnameFromURL: %s", err)
	}

	return parsed.Host, nil
}
