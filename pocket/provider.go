package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	contentTypeJSON      = "application/json; charset=UTF-8"
	urlPathGetNode       = "v1/query/node"
	urlPathGetHeight     = "v1/query/height"
	urlPathSimulateRelay = "v1/client/sim"
)

type Provider interface {
	Height() (uint, error)
	Servicer(address string) (Node, error)
	SimulateRelay(chainID, path string, payload json.RawMessage) (json.RawMessage, error)
}

func NewProvider(c http.Client, pocketRpcURL string) Provider {
	return provider{
		client:       c,
		pocketRpcURL: pocketRpcURL,
	}
}

type provider struct {
	client       http.Client
	pocketRpcURL string
}

func (p provider) Height() (uint, error) {
	return 0, nil
}

func (p provider) Servicer(address string) (Node, error) {
	return Node{}, nil
}

func (p provider) SimulateRelay(chainID, path string, payload json.RawMessage) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/%s", p.pocketRpcURL, urlPathSimulateRelay)

	simRequest := relayRequest{
		RelayNetworkID: chainID,
		Payload: relayRequestPayload{
			Data:    string(payload),
			Method:  "POST",
			Path:    path,
			Headers: make(map[string]string, 0),
		},
	}

	resp, err := p.doRequest(url, simRequest)
	if err != nil {
		return nil, fmt.Errorf("provider.SimulateRelay: %s", err)
	}

	return resp, nil
}

func (p provider) doRequest(url string, reqObj interface{}) ([]byte, error) {
	var reqBody []byte
	var err error
	if reqObj != nil {
		reqBody, err = json.Marshal(reqObj)
		if err != nil {
			return nil, fmt.Errorf("doRequest: %s", err)
		}
	}
	req := bytes.NewBuffer(reqBody)

	clientReq, err := http.NewRequest(http.MethodPost, url, req)
	if err != nil {
		return nil, fmt.Errorf("doRequest: %s", err)
	}
	clientReq.Header.Set("Content-type", contentTypeJSON)

	resp, err := p.client.Do(clientReq)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("doRequest: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("provider.doRequest: got unexpected response status %s - %s", resp.Status, string(reqBody)))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doRequest: %s", err)
	}

	return body, nil
}
