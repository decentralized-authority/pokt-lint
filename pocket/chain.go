package pocket

// Chain represents a pocket network chain as defined by https://docs.pokt.network/home/resources/references/supported-blockchains
type Chain struct {
	ID   string
	Name string
}

var allChains = map[string]string{
	// revenue earning
	ChainIDAlgorand:            "Algorand",
	ChainIDAVAX:                "AVAX",
	ChainIDBSC:                 "BSC",
	ChainIDBSCArchival:         "BSC Archival",
	ChainIDBoba:                "Boba",
	ChainIDDeFiKingdoms:        "DeFi Kingdoms",
	ChainIDSwimmer:             "Swimmer",
	ChainIDETH:                 "ETH",
	ChainIDETHArchival:         "ETH Archival",
	ChainIDAVAXArchival:        "AVAX Archival",
	ChainIDETHArchivalTrace:    "ETH Archival Trace",
	ChainIDETHGoerli:           "ETH Goerli",
	ChainIDETHKovan:            "ETH Kovan",
	ChainIDETHRinkeby:          "ETH Rinkeby",
	ChainIDETHRopsten:          "ETH Ropsten",
	ChainIDFantom:              "Fantom",
	ChainIDFUSE:                "FUSE",
	ChainIDFUSEArchival:        "FUSE Archival",
	ChainIDGnosisChain:         "Gnosis Chain",
	ChainIDGnosisChainArchival: "Gnosis Chain Archival",
	ChainIDHarmonyShard0:       "HMY 0",
	ChainIDIoTeX:               "IoTeX",
	ChainIDOKExChain:           "OKExChain",
	ChainIDPOKT:                "POKT",
	ChainIDPolygon:             "Polygon",
	ChainIDPolygonArchival:     "Polygon Archival",
	ChainIDSolana:              "Solana",

	// non-revenue
	ChainIDAlgorandArchival:        "Algorand Archival",
	ChainIDAlgorandTestnet:         "Algorand Testnet",
	ChainIDAlgorandTestnetArchival: "Algorand Testnet Archival",
	ChainIDArweave:                 "Arweave",
	ChainIDAVAXFuji:                "AVAX Fuji",
	ChainIDBSCTestnet:              "BSC Testnet",
	ChainIDBSCTestnetArchival:      "BSC Testnet Archival",
	ChainIDEvmos:                   "Evmos",
	ChainIDHarmonyShard0Archival:   "Harmony Shard 0 Archival",
	ChainIDHarmonyShard1:           "Harmony Shard 1",
	ChainIDHarmonyShard1Archival:   "Harmony Shard 1 Archival",
	ChainIDHarmonyShard2:           "Harmony Shard 2",
	ChainIDHarmonyShard2Archival:   "Harmony Shard 2 Archival",
	ChainIDHarmonyShard3:           "Harmony Shard 3",
	ChainIDHarmonyShard3Archival:   "Harmony Shard 3 Archival",
	ChainIDMoonbeam:                "Moonbeam",
	ChainIDMoonriver:               "Moonriver",
	ChainIDNEAR:                    "NEAR",
	ChainIDPolygonMumbai:           "Polygon Mumbai",
	ChainIDPolygonMumbaiArchival:   "Polygon Mumbai Archival",
	ChainIDSolanaTestnet:           "Solana Testnet",
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
