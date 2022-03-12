package rpc

import "net/http"

const (
	pocketQueryHeight  = `{}`
	btcGetBlockCount   = `{"jsonrpc":"1.0","id":"curltest","method":"getblockcount","params":[]}`
	avaxIsBootstrapped = `{"jsonrpc":"2.0","id":1,"method":"info.isBootstrapped","params":{"chain":"X"}}`
	ethBlockNumber     = `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":64}`
	hmyBlockNumber     = `{"jsonrpc":"2.0","method":"hmy_blockNumber","params":[],"id":64}`
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
		return Payload{Method: http.MethodPost, Path: "/v1/query/height", Data: pocketQueryHeight}
	case "0002":
		return Payload{Method: http.MethodPost, Path: "/", Data: btcGetBlockCount}
	case "0003":
		return Payload{Method: http.MethodPost, Path: "/ext/info", Data: avaxIsBootstrapped}
	case "0040":
		return Payload{Method: http.MethodPost, Path: "/", Data: hmyBlockNumber}
	default:
		return Payload{Method: http.MethodPost, Path: "/", Data: ethBlockNumber}
	}
}
