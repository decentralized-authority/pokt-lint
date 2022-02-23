package pocket

import "time"

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

type Session struct {
	ChainID      string
	Height       uint
	AppPublicKey string
	NumRelays    uint
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

type chainResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type balanceRequest struct {
	Address string `json:"address"`
}

type balanceResponse struct {
	Balance uint `json:"balance"`
}
