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

// NewService returns a new relaying.Service
func NewService(nodeID, nodeAddress string, chains []string, httpClient http.Client) (Service, error) {
	var err error
	chainObjects := make([]pocket.Chain, len(chains))
	pocketProvider := pocket.NewProvider(nodeAddress, httpClient)

	for i, c := range chains {
		if chainObjects[i], err = pocket.ChainFromID(c); err != nil {
			return service{}, fmt.Errorf("relaying.NewService: %s", err)
		}
	}

	nc := service{
		pocketProvider: pocketProvider,
		nodeID:         nodeID,
		nodeURL:        nodeAddress,
		nodeChains:     chainObjects,
	}

	if err := nc.init(); err != nil {
		return service{}, fmt.Errorf("relaying.NewService: %s", err)
	}

	return nc, nil
}

type errResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type service struct {
	pocketProvider pocket.Provider
	nodeID         string
	nodeURL        string
	nodeChains     []pocket.Chain
}

func (s service) RunRelayTests(_ context.Context, numSamples int64) (map[string]RelayTestResult, error) {
	if len(s.nodeChains) < 1 {
		return nil, fmt.Errorf("no chains for node %s", s.nodeID)
	}

	simIsEnabled, err := s.pocketProvider.SimulateRelayIsEnabled()
	if err != nil {
		return nil, fmt.Errorf("RunRelayTests: %s", err)
	}
	if !simIsEnabled {
		return nil, fmt.Errorf("simulateRelay is not enabled")
	}

	chains := make(map[string]RelayTestResult, len(s.nodeChains))
	for _, chain := range s.nodeChains {
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
			res, err := s.pocketProvider.SimulateRelay(req)
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

func (s *service) init() error {
	if len(s.nodeChains) > 0 {
		return nil
	}

	node, err := s.pocketProvider.Servicer(s.nodeID)
	if err != nil {
		return fmt.Errorf("init: %s: %s", s.nodeID, err)
	}

	s.nodeURL = node.ServiceURL
	s.nodeChains = node.Chains
	return nil
}

// NewNodeChecker returns a node checker relaying service
func NewNodeChecker(nodeID, nodeAddress string, chains []string, provider pocket.Provider) (Service, error) {
	var err error
	empty := service{}
	chainObjects := make([]pocket.Chain, len(chains))

	for i, c := range chains {
		if chainObjects[i], err = pocket.ChainFromID(c); err != nil {
			return service{}, fmt.Errorf("relaying.NewNodeChecker: %s", err)
		}
	}

	nc := service{
		pocketProvider: provider,
		nodeID:         nodeID,
		nodeURL:        nodeAddress,
		nodeChains:     chainObjects,
	}

	if err := nc.init(); err != nil {
		return empty, fmt.Errorf("relaying.NewNodeChecker: %s", err)
	}

	return nc, nil
}
