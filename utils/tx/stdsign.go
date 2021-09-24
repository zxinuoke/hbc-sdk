package tx

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/zxinuoke/hbc-sdk/utils"
)

// StdSignDoc def
type StdSignDoc struct {
	//AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	ChainID  string            `json:"chain_id" yaml:"chain_id"`
	Fee      json.RawMessage   `json:"fee" yaml:"fee"`
	Memo     string            `json:"memo" yaml:"memo"`
	Msgs     []json.RawMessage `json:"msgs" yaml:"msgs"`
	Sequence uint64            `json:"sequence" yaml:"sequence"`
}

// StdSignMsg def
type StdSignMsg struct {
	ChainID string `json:"chain_id"`
	//AccountNumber int64       `json:"account_number"`
	Sequence int64       `json:"sequence"`
	Msgs     []utils.Msg `json:"msgs"`
	Memo     string      `json:"memo"`
	Fee      StdFee      `json:"fee" yaml:"fee"`
}

// StdSignature Signature
type StdSignature struct {
	crypto.PubKey `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
	//AccountNumber int64  `json:"account_number"`
	//Sequence      int64  `json:"sequence"`
}

// Bytes gets message bytes
func (msg StdSignMsg) Bytes() []byte {
	return StdSignBytes(msg.ChainID, msg.Sequence, msg.Msgs, msg.Memo, msg.Fee)
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, sequence int64, msgs []utils.Msg, memo string, fee StdFee) []byte {
	msgsBytes := make([]json.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}

	bz, err := codec.Cdc.MarshalJSON(StdSignDoc{
		//AccountNumber: uint64(accnum),
		ChainID:  chainID,
		Fee:      json.RawMessage(fee.Bytes()),
		Memo:     memo,
		Msgs:     msgsBytes,
		Sequence: uint64(sequence),
	})
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(bz)
}
