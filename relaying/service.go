package relaying

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"github.com/itsnoproblem/pokt-lint/rpc"
	"github.com/itsnoproblem/pokt-lint/timer"
)

// Service represents a relaying service
type Service interface {
	RunRelayTests(ctx context.Context) (map[string]RelayTestResult, error)
}

type nodeChecker struct {
	pocketProvider pocket.Provider
	nodeID         string
	nodeURL        string
	nodeChains     []pocket.Chain
}

func (c nodeChecker) RunRelayTests(_ context.Context) (map[string]RelayTestResult, error) {
	if len(c.nodeChains) < 1 {
		return nil, fmt.Errorf("no chains for node %s", c.nodeID)
	}

	chains := make(map[string]RelayTestResult, len(c.nodeChains))
	for _, chain := range c.nodeChains {
		var success bool
		var msg string

		req := pocket.RelayRequest{
			RelayNetworkID: chain.ID,
			Payload:        rpc.NewPayload(chain.ID),
		}

		t := timer.Start()
		res, err := c.pocketProvider.SimulateRelay(req)
		if err != nil {
			success = false
			msg = err.Error()
		} else {
			success = true
			msg = "OK"
		}

		chains[chain.ID] = RelayTestResult{
			ChainID:       chain.ID,
			Successful:    success,
			Message:       msg,
			StatusCode:    res.StatusCode,
			DurationMS:    float64(t.Elapsed().Microseconds()) / 1000,
			RelayRequest:  req.Payload,
			RelayResponse: res,
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

// NewNodeChecker returns a node checker relaying service
func NewNodeChecker(nodeID, nodeAddress string, chains []string, httpClient http.Client) (Service, error) {
	var err error
	empty := nodeChecker{}
	chainObjects := make([]pocket.Chain, len(chains))
	pocketProvider := pocket.NewProvider(httpClient, nodeAddress)

	for i, c := range chains {
		if chainObjects[i], err = pocket.ChainFromID(c); err != nil {
			return nodeChecker{}, fmt.Errorf("relaying.NewNodeChecker: %s", err)
		}
	}

	nc := nodeChecker{
		pocketProvider: pocketProvider,
		nodeID:         nodeID,
		nodeURL:        nodeAddress,
		nodeChains:     chainObjects,
	}

	if err := nc.init(); err != nil {
		return empty, fmt.Errorf("relaying.NewNodeChecker: %s", err)
	}

	return nc, nil
}
