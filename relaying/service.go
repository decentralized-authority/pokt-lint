package relaying

import (
	"context"
	"errors"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"github.com/itsnoproblem/pokt-lint/rpc"
	"github.com/itsnoproblem/pokt-lint/timer"
	nethttp "net/http"
)

type Service interface {
	RunRelayTests(ctx context.Context) (map[string]RelayTestResult, error)
}

type nodeChecker struct {
	pocketProvider pocket.Provider
	nodeID         string
	nodeURL        string
	nodeChains     []pocket.Chain
}

type RelayTestResult struct {
	ChainID    string                 `json:"chain_id"`
	Successful bool                   `json:"success"`
	Data       map[string]interface{} `json:"data"`
	DurationMS float64                `json:"duration_ms"`
}

func NewNodeChecker(nodeID, nodeAddress string, chains []string, httpClient nethttp.Client) (*nodeChecker, error) {
	var err error
	empty := nodeChecker{}
	chainObjects := make([]pocket.Chain, len(chains))
	pocketProvider := pocket.NewProvider(httpClient, nodeAddress)

	for i, c := range chains {
		if chainObjects[i], err = pocket.ChainFromID(c); err != nil {
		}
	}

	nc := nodeChecker{
		pocketProvider: pocketProvider,
		nodeID:         nodeID,
		nodeURL:        nodeAddress,
		nodeChains:     chainObjects,
	}

	if err := nc.init(); err != nil {
		return &empty, fmt.Errorf("relaying.NewNodeChecker: %s", err)
	}

	return &nc, nil
}

func (c *nodeChecker) RunRelayTests() (map[string]RelayTestResult, error) {
	if len(c.nodeChains) < 1 {
		return nil, errors.New(fmt.Sprintf("No chains for node %s", c.nodeID))
	}

	chains := make(map[string]RelayTestResult, len(c.nodeChains))
	for _, chain := range c.nodeChains {
		var success bool
		msg := make(map[string]interface{})

		req := pocket.RelayRequest{
			RelayNetworkID: chain.ID,
			Payload:        rpc.NewPayload(chain.ID),
		}

		t := timer.Start()
		res, err := c.pocketProvider.SimulateRelay(req)
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

		chains[chain.ID] = RelayTestResult{
			ChainID:    chain.ID,
			Successful: success,
			Data:       msg,
			DurationMS: float64(t.Elapsed().Microseconds()) / 1000,
		}
	}
	return chains, nil
}

func (c *nodeChecker) init() error {
	if len(c.nodeChains) > 0 {
		return nil
	}

	node, err := c.pocketProvider.Servicer(c.nodeID)
	if err != nil {
		return fmt.Errorf("init: %s", err)
	}

	c.nodeURL = node.ServiceURL
	c.nodeChains = node.Chains
	return nil
}
