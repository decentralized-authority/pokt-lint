package pocket

// Chain represents a pocket network chain as defined by https://docs.pokt.network/home/resources/references/supported-blockchains
type Chain struct {
	ID   string
	Name string
}

var allChains = map[string]string{
	// revenue earning
	"0029": "Algorand",
	"0003": "AVAX",
	"0004": "BSC",
	"0010": "BSC Archival",
	"0048": "Boba",
	"03DF": "DeFi Kingdoms",
	"0021": "ETH",
	"0022": "ETH Archival",
	"00A3": "AVAX Archival",
	"0028": "ETH Archival Trace",
	"0026": "ETH Goerli",
	"0024": "ETH Kovan",
	"0025": "ETH Rinkeby",
	"0023": "ETH Ropsten",
	"0049": "Fantom",
	"0005": "FUSE",
	"000A": "FUSE Archival",
	"0027": "Gnosis Chain",
	"000C": "Gnosis Chain Archival",
	"0040": "HMY 0",
	"0044": "IoTeX",
	"0047": "OKExChain",
	"0001": "POKT",
	"0009": "Polygon",
	"000B": "Polygon Archival",
	"0006": "Solana",

	// non-revenue
	"000D": "Algorand Archival",
	"0045": "Algorand Testnet",
	"0A45": "Algorand Testnet Archival",
	"0030": "Arweave",
	"000E": "AVAX Fuji",
	"0011": "BSC Testnet",
	"0012": "BSC Testnet Archival",
	"0046": "Evmos",
	"0A40": "Harmony Shard 0 Archival",
	"0041": "Harmony Shard 1",
	"0A41": "Harmony Shard 1 Archival",
	"0042": "Harmony Shard 2",
	"0A42": "Harmony Shard 2 Archival",
	"0043": "Harmony Shard 3",
	"0A43": "Harmony Shard 3 Archival",
	"0050": "Moonbeam",
	"0051": "Moonriver",
	"0052": "NEAR",
	"000F": "Polygon Mumbai",
	"00AF": "Polygon Mumbai Archival",
	"0031": "Solana Testnet",
}

// ChainFromID returns the chain for a given ID, or an error if not found
func ChainFromID(id string) (Chain, error) {
	name, ok := allChains[id]
	if !ok {
		name = "Unknown"
	}

	return Chain{
		ID:   id,
		Name: name,
	}, nil
}
