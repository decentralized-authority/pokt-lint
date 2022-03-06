package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	gohttp "net/http"
	"strconv"

	"github.com/itsnoproblem/pokt-lint/http"
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
	SimulateRelay(req RelayRequest) (RelayResponse, error)
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

	body, _, err := p.doRequest(url, nodeRequest)
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

func bytesToMap(b []byte) map[string]interface{} {
	thing := make(map[string]interface{})
	_ = json.Unmarshal(b, &thing)
	return thing
}

func (p provider) SimulateRelay(simRequest RelayRequest) (RelayResponse, error) {
	url := fmt.Sprintf("%s/%s", p.pocketRpcURL, urlPathSimulateRelay)

	respBody, statusCode, err := p.doRequest(url, simRequest)
	if err != nil {
		return RelayResponse{
			StatusCode: statusCode,
			Data:       respBody,
		}, fmt.Errorf("pocketProvider.SimulateRelay: %s", err)
	}

	return RelayResponse{
		StatusCode: statusCode,
		Data:       respBody,
	}, nil
}

func (p provider) doRequest(url string, reqObj interface{}) ([]byte, int, error) {
	var reqBody []byte
	var err error
	if reqObj != nil {
		reqBody, err = json.Marshal(reqObj)
		if err != nil {
			return nil, 500, fmt.Errorf("doRequest got error encoding request: %s", err)
		}
	}
	req := bytes.NewBuffer(reqBody)

	clientReq, err := gohttp.NewRequest(gohttp.MethodPost, url, req)
	if err != nil {
		return nil, 500, fmt.Errorf("doRequest got error creating request: %s", err)
	}
	clientReq.Header.Set("Content-type", contentTypeJSON)

	resp, err := p.client.Do(clientReq)
	body := make([]byte, 0)
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Default().Printf("doRequest got error closing response body: %s", err)
			}
		}(resp.Body)

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, fmt.Errorf("doRequest got error reading response body: %s", err)
		}

	}

	if err != nil {
		return nil, 500, fmt.Errorf("doRequest: %s", err)
	}

	return body, resp.StatusCode, nil
}
