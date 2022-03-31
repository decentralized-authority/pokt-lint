package relaying

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/http"
	"github.com/itsnoproblem/pokt-lint/maths"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"github.com/itsnoproblem/pokt-lint/rpc"
	"github.com/itsnoproblem/pokt-lint/timer"
)

// Service represents a relaying service
type Service interface {
	RunRelayTests(ctx context.Context, numSamples int64) (map[string]RelayTestResult, error)
}

type errResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type nodeChecker struct {
	pocketProvider pocket.Provider
	nodeID         string
	nodeURL        string
	nodeChains     []pocket.Chain
}

func (c nodeChecker) RunRelayTests(_ context.Context, numSamples int64) (map[string]RelayTestResult, error) {
	if len(c.nodeChains) < 1 {
		return nil, fmt.Errorf("no chains for node %s", c.nodeID)
	}

	chains := make(map[string]RelayTestResult, len(c.nodeChains))
	for _, chain := range c.nodeChains {
		req := pocket.RelayRequest{
			RelayNetworkID: chain.ID,
			Payload:        rpc.NewPayload(chain.ID),
		}

		result := RelayTestResult{
			ChainID:        chain.ID,
			ChainName:      chain.Name,
			RelayRequest:   req.Payload,
			RelayResponses: make([]RelayTestSample, numSamples),
		}

		allDurations := make([]float64, numSamples)

		for i := int64(0); i < numSamples; i++ {
			t := timer.Start()
			res, err := c.pocketProvider.SimulateRelay(req)
			result.StatusCode = res.StatusCode

			if err != nil {
				result.Successful = false
				result.Message = err.Error()

			} else if res.StatusCode != 200 {
				var relayErr errResponse
				_ = json.Unmarshal(res.Data, &relayErr)
				result.Successful = false
				result.Message = relayErr.Message
				result.StatusCode = relayErr.Code

			} else {
				result.Successful = true
				result.Message = "OK"
			}

			duration := t.Elapsed().Microseconds()
			durationMS := float64(duration) / 1000
			allDurations[i] = durationMS

			result.RelayResponses[i] = RelayTestSample{
				DurationMS: durationMS,
				StatusCode: res.StatusCode,
				Data:       res.Data,
			}
		}

		result.DurationAvgMS = maths.Mean(allDurations)
		result.DurationMedianMS = maths.Median(allDurations)
		result.DurationMaxMS = maths.Max(allDurations)
		result.DurationMinMS = maths.Min(allDurations)
		chains[chain.ID] = result
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
