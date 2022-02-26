package rpc

import "net/http"

const (
	PocketQueryHeight  = `{}`
	BtcGetBlockCount   = `{"jsonrpc":"1.0","id":"curltest","method":"getblockcount","params":[]}`
	AvaxIsBootstrapped = `{"jsonrpc":"2.0","id":1,"method":"info.isBootstrapped","params":{"chain":"X"}}`
	EthBlockNumber     = `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":64}`
	HmyBlockNumber     = `{"jsonrpc":"2.0","method":"hmy_blockNumber","params":[],"id":64}`
)

type Payload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

func NewPayload(chainID string) Payload {
	switch chainID {
	case "0001":
		return Payload{Method: http.MethodPost, Path: "/v1/query/height", Data: PocketQueryHeight}
	case "0002":
		return Payload{Method: http.MethodPost, Path: "/", Data: BtcGetBlockCount}
	case "0003":
		return Payload{Method: http.MethodPost, Path: "/ext/info", Data: AvaxIsBootstrapped}
	case "0040":
		return Payload{Method: http.MethodPost, Path: "/", Data: HmyBlockNumber}
	default:
		return Payload{Method: http.MethodPost, Path: "/", Data: EthBlockNumber}
	}
}
