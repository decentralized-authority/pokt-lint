package pocket

import (
	"fmt"
	"github.com/itsnoproblem/pokt-lint/rpc"
)

type RelayRequest struct {
	RelayNetworkID string      `json:"relay_network_id"`
	Payload        rpc.Payload `json:"payload"`
}

type RelayError struct {
	Code int
	Err  error
}

func (r RelayError) Error() string {
	return fmt.Sprintf("(%d) %s", r.Code, r.Err)
}

func NewRelayError(code int, err error) RelayError {
	return RelayError{
		Code: code,
		Err:  err,
	}
}
