package relay

type Request struct {
	RelayNetworkID string         `json:"relay_network_id"`
	Payload        RequestPayload `json:"payload"`
}

func (r *Request) Path() string {
	switch r.RelayNetworkID {
	case "0003":
		return "/ext/info"
	case "0001":
		return "/v1/query/height"
	}

	return ""
}

type RequestPayload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}
