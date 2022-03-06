package pocket

import (
	"encoding/json"
	"github.com/itsnoproblem/pokt-lint/rpc"
)

type RelayRequest struct {
	RelayNetworkID string      `json:"relay_network_id"`
	Payload        rpc.Payload `json:"payload"`
}

type RelayResponse struct {
	StatusCode int             `json:"status_code"`
	Data       json.RawMessage `json:"data"`
}
