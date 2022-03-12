package pocket

import (
	"encoding/json"
	"github.com/itsnoproblem/pokt-lint/rpc"
)

// RelayRequest is the format of a simulateRelay request
type RelayRequest struct {
	RelayNetworkID string      `json:"relay_network_id"`
	Payload        rpc.Payload `json:"payload"`
}

// RelayResponse is the format of the response expected from a call to simulateRelay
type RelayResponse struct {
	StatusCode int             `json:"status_code"`
	Data       json.RawMessage `json:"data"`
}
