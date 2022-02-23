package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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
	SimulateRelay(chainID, path string, payload json.RawMessage) (map[string]interface{}, error)
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
	var fail = func(err error) (Node, error) {
		return Node{}, fmt.Errorf("Services: %s", err)
	}

	url := fmt.Sprintf("%s/%s", p.pocketRpcURL, urlPathGetNode)
	nodeRequest := queryNodeRequest{Address: address}
	var nodeResponse queryNodeResponse

	body, err := p.doRequest(url, nodeRequest)
	if err != nil {
		return fail(err)
	}

	err = json.Unmarshal(body, &nodeResponse)
	if err != nil {
		return fail(err)
	}

	chains := make([]Chain, len(nodeResponse.Chains))
	for i, chainID := range nodeResponse.Chains {
		ch, err := ChainFromID(chainID)
		if err != nil {
			fail(err)
		}

		chains[i] = ch
	}

	stakedBal, err := strconv.ParseUint(nodeResponse.StakedBalance, 10, 64)
	if err != nil {
		return Node{}, fmt.Errorf("Node: %s", err)
	}

	return Node{
		Address:       nodeResponse.Address,
		Pubkey:        nodeResponse.Pubkey,
		ServiceURL:    nodeResponse.ServiceURL,
		StakedBalance: uint(stakedBal),
		IsJailed:      nodeResponse.IsJailed,
		Chains:        chains,
		IsSynced:      false,
	}, nil
}

func (p provider) SimulateRelay(chainID, path string, payload json.RawMessage) (map[string]interface{}, error) {
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
		return nil, err
	}

	s, _ := strconv.Unquote(string(resp))
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (p provider) doRequest(url string, reqObj interface{}) ([]byte, error) {
	var reqBody []byte
	var err error
	if reqObj != nil {
		reqBody, err = json.Marshal(reqObj)
		if err != nil {
			return nil, NewRelayError(500, err)
		}
	}
	req := bytes.NewBuffer(reqBody)

	clientReq, err := http.NewRequest(http.MethodPost, url, req)
	if err != nil {
		return nil, NewRelayError(500, err)
	}
	clientReq.Header.Set("Content-type", contentTypeJSON)

	resp, err := p.client.Do(clientReq)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Default().Printf("error closing response body: %s", err)
		}
	}(resp.Body)
	if err != nil {
		return nil, NewRelayError(500, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewRelayError(resp.StatusCode, err)
	}

	log.Default().Printf("Pocket Provider: (%d) %s", resp.StatusCode, url)
	if resp.StatusCode != http.StatusOK {
		var str string
		_ = json.Unmarshal(body, &str)
		return nil, NewRelayError(resp.StatusCode, errors.New(str))
	}

	return body, nil
}
