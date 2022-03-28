package pocket

import (
	"fmt"
)

// Chain represents a pocket network chain as defined by https://docs.pokt.network/home/resources/references/supported-blockchains
type Chain struct {
	ID   string
	Name string
}

var allChains = map[string]string{
	"0003": "AVAX",
	"00A3": "AVAX Archival",
	"0004": "BSC",
	"0010": "BSC Archival",
	"0021": "ETH",
	"0022": "ETH Archival",
	"0028": "ETH Archival Trace",
	"0026": "ETH Goerli",
	"0024": "ETH Kovan",
	"0025": "ETH Rinkeby",
	"0023": "ETH Ropsten",
	"0005": "FUSE",
	"000A": "FUSE Archival",
	"0040": "HMY 0",
	"0044": "IoTeX",
	"0001": "POKT",
	"0002": "POKT testnet",
	"0009": "Polygon",
	"000B": "Polygon Archival",
	"000F": "Polygon Mumbai",
	"0006": "Solana",
	"0027": "Gnosis Chain",
	"000C": "Gnosis Chain Archival",
	"0029": "Algorand",
	"0047": "OKExChain",
}

// ChainFromID returns the chain for a given ID, or an error if not found
func ChainFromID(id string) (Chain, error) {
	name, ok := allChains[id]
	if !ok {
		return Chain{}, fmt.Errorf("ChainFromID: unknown chain %s", id)
	}

	return Chain{
		ID:   id,
		Name: name,
	}, nil
}
