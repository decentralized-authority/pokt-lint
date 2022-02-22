package linting

import (
	"encoding/json"
	"fmt"
)

type RelayTestRequest struct {
	ChainID string
	Path    string
	Payload json.RawMessage
}

func RPCRequestForChainID(chainID string) (RelayTestRequest, error) {
	emptySlice := make([]string, 0)
	payload := make(map[string]interface{})
	var path = ""
	var err error

	switch chainID {
	case "0001":
		path = "/v1/query/height"

	case "0003":
		path = "/ext/info"
		payload = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "info.getNetworkID",
		}

	case "0040":

		payload = map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "hmyv2_getEpoch",
			"params":  emptySlice,
			"id":      1,
		}

	default:
		path = "/"
		payload = map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_chainId",
			"params":  emptySlice,
			"id":      1,
		}
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return RelayTestRequest{}, fmt.Errorf("RPCRequestForChainID: %s", err)
	}

	return RelayTestRequest{
		ChainID: chainID,
		Path:    path,
		Payload: encoded,
	}, nil
}

