package pocket

import "time"

// Node represents a pocket network servicer node
type Node struct {
	Address           string
	Pubkey            string
	Balance           uint
	StakedBalance     uint
	ServiceURL        string
	IsJailed          bool
	Chains            []Chain
	IsSynced          bool
	LatestBlockHeight uint
	LatestBlockTime   time.Time
}

type queryNodeRequest struct {
	Address string `json:"address"`
}

type queryNodeResponse struct {
	Address       string   `json:"address"`
	Pubkey        string   `json:"public_key"`
	Chains        []string `json:"chains"`
	IsJailed      bool     `json:"jailed"`
	ServiceURL    string   `json:"service_url"`
	StakedBalance string   `json:"tokens"`
}
