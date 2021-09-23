package hbc

import (
	"encoding/json"
)

type HbcGas struct {
	Fee string `json:"fee"`
	Gas string `json:"gas"`
}

type BaseResponse struct {
	Error string `json:"error"`
}

type BlockData struct {
	BaseResponse
	Block struct {
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`

		Header struct {
			AppHash            string      `json:"app_hash"`
			ChainID            string      `json:"chain_id"`
			ConsensusHash      string      `json:"consensus_hash"`
			DataHash           string      `json:"data_hash"`
			EvidenceHash       string      `json:"evidence_hash"`
			Height             json.Number `json:"height"`
			LastCommitHash     string      `json:"last_commit_hash"`
			LastResultsHash    string      `json:"last_results_hash"`
			NextValidatorsHash string      `json:"next_validators_hash"`
			NumTxs             string      `json:"num_txs"`
			ProposerAddress    string      `json:"proposer_address"`
			Time               string      `json:"time"`
			TotalTxs           string      `json:"total_txs"`
			ValidatorsHash     string      `json:"validators_hash"`
		} `json:"header"`
	} `json:"block"`
}

type TxAmount struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type TxMsg struct {
	Type  string `json:"type"`
	Value struct {
		Amount      []TxAmount `json:"amount"`
		FromAddress string     `json:"from_address"`
		ToAddress   string     `json:"to_address"`
	} `json:"value"`
}

type AddressData struct {
	BaseResponse
	Height json.Number `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			AccountNumber json.Number `json:"account_number"`
			Address       string      `json:"address"`
			Coins         []struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"coins"`
			Sequence json.Number `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}

type TxData struct {
	BaseResponse
	GasUsed   json.Number `json:"gas_used"`
	GasWanted json.Number `json:"gas_wanted"`
	Height    json.Number `json:"height"`
	Txhash    string      `json:"txhash"`
	Timestamp string      `json:"timestamp"`
	Logs      []struct {
		Log      string `json:"log"`
		MsgIndex int64  `json:"msg_index"`
		Success  bool   `json:"success"`
	} `json:"logs"`

	Tx struct {
		Type  string `json:"type"`
		Value struct {
			Fee struct {
				Amount []TxAmount  `json:"amount"`
				Gas    json.Number `json:"gas"`
			} `json:"fee"`
			Memo string  `json:"memo"`
			Msg  []TxMsg `json:"msg"`
		} `json:"value"`
	} `json:"tx"`
}

type TxResponse struct {
	Height    json.Number `json:"height"`
	TxHash    string      `json:"txhash"`
	Codespace string      `json:"codespace,omitempty"`
	Code      json.Number `json:"code,omitempty"`
	Data      string      `json:"data,omitempty"`
	RawLog    string      `json:"raw_log,omitempty"`
	Info      string      `json:"info,omitempty"`
	GasWanted json.Number `json:"gas_wanted,omitempty"`
	GasUsed   json.Number `json:"gas_used,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

type AssetData struct {
	BaseResponse
	Height json.Number `json:"height"`
	Result struct {
		Available []struct {
			Amount string `json:"amount"`
			Denom  string `json:"denom"`
		} `json:"available"`
	} `json:"result"`
}
