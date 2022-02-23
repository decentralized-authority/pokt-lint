package pocket

import (
	"fmt"
)

type relayRequest struct {
	RelayNetworkID string              `json:"relay_network_id"`
	Payload        relayRequestPayload `json:"payload"`
}

type relayRequestPayload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

type RelayError struct {
	Code int
	Err  error
}

func (r RelayError) Error() string {
	return fmt.Sprintf("(%s) %s", r.Code, r.Err)
}

func NewRelayError(code int, err error) RelayError {
	return RelayError{
		Code: code,
		Err:  err,
	}
}
