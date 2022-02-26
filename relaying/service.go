package relaying

import (
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pocket"
	"github.com/itsnoproblem/pokt-lint/rpc"
)

type Service interface {
	SimulateRelay(servicerUrl, chain, path string, payload rpc.Payload) (map[string]interface{}, error)
}

func NewService(p pocket.Provider) Service {
	return service{
		pocketProvider: p,
	}
}

type service struct {
	pocketProvider pocket.Provider
}

func (s service) SimulateRelay(servicerUrl, chainID, path string, rpcPayload rpc.Payload) (map[string]interface{}, error) {
	resp, err := s.pocketProvider.SimulateRelay(servicerUrl, chainID, rpcPayload)
	if err != nil {
		return nil, fmt.Errorf("relaying.SimulateRelay: %s", err)
	}

	return resp, nil
}
