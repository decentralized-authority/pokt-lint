package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	nethttp "net/http"
	"strconv"
	"strings"

	"github.com/itsnoproblem/pokt-lint/http"
)

const (
	contentTypeJSON      = "application/json; charset=UTF-8"
	urlPathGetNode       = "v1/query/node"
	urlPathSimulateRelay = "v1/client/sim"
)

// Provider is a pocket network client
type Provider interface {
	Height() (uint, error)
	Servicer(address string) (Node, error)
	SimulateRelay(req RelayRequest) (RelayResponse, error)
	SimulateRelayIsEnabled() (bool, error)
}

//type HTTPClient interface {
//	Do(req *nethttp.Request) (*nethttp.Response, error)
//	Get(url string) (*nethttp.Response, error)
//	Options(url string) (*nethttp.Response, error)
//}

// NewProvider returns a new pocket provider
func NewProvider(pocketURL string, client http.Client) Provider {
	return provider{
		client:    client,
		pocketURL: pocketURL,
	}
}

type provider struct {
	client    http.Client
	pocketURL string
}

func (p provider) Height() (uint, error) {
	return 0, nil
}

func (p provider) Servicer(address string) (Node, error) {
	var fail = func(err error) (Node, error) {
		return Node{}, fmt.Errorf("provider.Servicer: %s", err)
	}

	url := fmt.Sprintf("%s/%s", p.pocketURL, urlPathGetNode)
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

	var stakedBal uint64
	if nodeResponse.StakedBalance != "" {
		stakedBal, err = strconv.ParseUint(nodeResponse.StakedBalance, 10, 64)
		if err != nil {
			return Node{}, fmt.Errorf("failed to parse staked balance (%s): %s", nodeResponse.StakedBalance, err)
		}
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

func (p provider) SimulateRelay(simRequest RelayRequest) (RelayResponse, error) {
	url := fmt.Sprintf("%s/%s", p.pocketURL, urlPathSimulateRelay)

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

func (p provider) SimulateRelayIsEnabled() (bool, error) {
	url := fmt.Sprintf("%s/%s", p.pocketURL, urlPathSimulateRelay)

	res, err := p.client.Options(url)
	if err != nil {
		return false, fmt.Errorf("SimulateRelayIsEnabled: %s", err)
	}

	if res.StatusCode != nethttp.StatusOK {
		return false, nil
	}

	return true, nil
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

	clientReq, err := nethttp.NewRequest(nethttp.MethodPost, url, req)
	if err != nil {
		return nil, 500, fmt.Errorf("doRequest got error creating request: %s", err)
	}
	clientReq.Header.Set("Content-type", contentTypeJSON)

	resp, err := p.client.Do(clientReq)
	if err != nil {
		return nil, 500, fmt.Errorf("doRequest: %s", err)
	}

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

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		body = []byte(fmt.Sprintf(`{ "body": "%s" }`, strings.Trim(string(body), "\n")))
		log.Default().Printf("formatting non json body: %s", body)
	}

	return body, resp.StatusCode, nil
}
