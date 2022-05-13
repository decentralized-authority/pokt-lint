package rpc

import "net/http"

const (
	emptyPayload       = `{}`
	btcGetBlockCount   = `{"jsonrpc":"1.0","id":"curltest","method":"getblockcount","params":[]}`
	avaxIsBootstrapped = `{"jsonrpc":"2.0","id":1,"method":"info.isBootstrapped","params":{"chain":"X"}}`
	avaxGetBlockchains = `{"jsonrpc":"2.0","id":1,"method":"platform.getBlockchains","params":{}}`
	ethBlockNumber     = `{"jsonrpc":"2.0","id":64,"method":"eth_blockNumber","params":[]}`
	hmyBlockNumber     = `{"jsonrpc":"2.0","id":1,"method":"hmyv2_getBlocks","params":[1,2,{"withSigners":false,"fullTx":false,"inclStaking":false}]}`
)

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
	case "0001":
		return Payload{Method: http.MethodPost, Path: "/v1/query/height", Data: emptyPayload}

	case "0002":
		return Payload{Method: http.MethodPost, Path: "/", Data: btcGetBlockCount}

	case "0003":
		return Payload{Method: http.MethodPost, Path: "/ext/info", Data: avaxIsBootstrapped}

	case "03DF":
		return Payload{Method: http.MethodPost, Path: "/ext/P", Data: avaxGetBlockchains}

	case "0029":
		return Payload{Method: http.MethodGet, Path: "/v2/status", Data: emptyPayload}

	case "0030":
		return Payload{Method: http.MethodGet, Path: "/info", Data: emptyPayload}

	case "0040":
	case "0A40":
	case "0041":
	case "0A41":
	case "0042":
	case "0A42":
	case "0043":
	case "0A43":
		return Payload{Method: http.MethodPost, Path: "/", Data: hmyBlockNumber}
	}

	// EVM compatible payload is default
	return Payload{Method: http.MethodPost, Path: "/", Data: ethBlockNumber}
}
