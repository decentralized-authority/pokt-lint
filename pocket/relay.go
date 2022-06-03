package pocket

import (
	"encoding/json"
	"net/http"
)

const (
	EmptyPayload       = `{}`
	AVAXIsBootstrapped = `{"jsonrpc":"2.0","id":1,"method":"info.isBootstrapped","params":{"chain":"X"}}`
	AVAXGetBlockchains = `{"jsonrpc":"2.0","id":1,"method":"platform.getBlockchains","params":{}}`
	ETHBlockNumber     = `{"jsonrpc":"2.0","id":64,"method":"eth_blockNumber","params":[]}`
	HMYBlockNumber     = `{"jsonrpc":"2.0","id":1,"method":"hmyv2_getBlocks","params":[1,2,{"withSigners":false,"fullTx":false,"inclStaking":false}]}`
)

// RelayRequest is the format of a simulateRelay request
type RelayRequest struct {
	RelayNetworkID string  `json:"relay_network_id"`
	Payload        Payload `json:"payload"`
}

// RelayResponse is the format of the response expected from a call to simulateRelay
type RelayResponse struct {
	StatusCode int             `json:"status_code"`
	Data       json.RawMessage `json:"data"`
}

// Payload represents an RPC payload sent to a pocket network servicer
type Payload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

// NewPayload returns a payload formatted for the specified chain
func NewPayload(chainID string) Payload {
	switch chainID {
	case ChainIDPOKT:
		return Payload{Method: http.MethodPost, Path: "/v1/query/height", Data: EmptyPayload}

	case ChainIDAVAX:
		return Payload{Method: http.MethodPost, Path: "/ext/info", Data: AVAXIsBootstrapped}

	case ChainIDDeFiKingdoms:
		return Payload{
			Method: http.MethodPost,
			Path:   "/ext/bc/q2aTwKuyzgs8pynF7UXBZCU7DejbZbZ6EUyHr3JQzYgwNPUPi/rpc",
			Data:   ETHBlockNumber,
		}

	case ChainIDSwimmer:
		return Payload{
			Method: http.MethodPost,
			Path:   "/ext/bc/2K33xS9AyP9oCDiHYKVrHe7F54h2La5D8erpTChaAhdzeSu2RX/rpc",
			Data:   ETHBlockNumber,
		}

	case ChainIDAlgorand:
		return Payload{Method: http.MethodGet, Path: "/v2/status", Data: EmptyPayload}

	case ChainIDArweave:
		return Payload{Method: http.MethodGet, Path: "/info", Data: EmptyPayload}

	case ChainIDHarmonyShard0:
	case ChainIDHarmonyShard0Archival:
	case ChainIDHarmonyShard1:
	case ChainIDHarmonyShard1Archival:
	case ChainIDHarmonyShard2:
	case ChainIDHarmonyShard2Archival:
	case ChainIDHarmonyShard3:
	case ChainIDHarmonyShard3Archival:
		return Payload{Method: http.MethodPost, Path: "/", Data: HMYBlockNumber}
	}

	// EVM compatible payload is default
	return Payload{Method: http.MethodPost, Path: "/", Data: ETHBlockNumber}
}
