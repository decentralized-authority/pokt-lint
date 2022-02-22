package relaying

import (
	"encoding/json"
	"fmt"
	"github.com/itsnoproblem/pokt-lint/pocket"
)

type Service interface {
	SimulateRelay(servicerUrl, chain, path string, payload map[string]interface{}) (json.RawMessage, error)
}

func NewService(p pocket.Provider) Service {
	return service{
		pocketProvider: p,
	}
}

type service struct {
	pocketProvider pocket.Provider
}

func (s service) SimulateRelay(servicerUrl, chainID, path string, payload map[string]interface{}) (json.RawMessage, error) {
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("relaying.SimulateRelay: %s", err)
	}

	resp, err := s.pocketProvider.SimulateRelay(servicerUrl, chainID, encodedPayload)
	if err != nil {
		return nil, fmt.Errorf("relaying.SimulateRelay: %s", err)
	}

	return resp, nil
}
