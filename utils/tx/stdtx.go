package tx

import (
	"fmt"

	"github.com/bluehelix-chain/hbc-sdk/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxGasWanted defines the max gas allowed.
const MaxGasWanted = uint64((1 << 63) - 1)

type StdFee struct {
	Amount sdk.Coins `json:"amount"`
	Gas    uint64    `json:"gas,omitempty"`
}

// NewStdFee returns a new instance of StdFee
func NewStdFee(gas uint64, amount sdk.Coins) StdFee {
	return StdFee{
		Amount: amount,
		Gas:    gas,
	}
}

// GetGas returns the fee's (wanted) gas.
func (fee StdFee) GetGas() uint64 {
	return fee.Gas
}

// GetAmount returns the fee's amount.
func (fee StdFee) GetAmount() sdk.Coins {
	return fee.Amount
}

// Bytes returns the encoded bytes of a StdFee.
func (fee StdFee) Bytes() []byte {
	if len(fee.Amount) == 0 {
		fee.Amount = sdk.NewCoins()
	}

	bz, err := codec.Cdc.MarshalJSON(fee)
	if err != nil {
		panic(err)
	}

	return bz
}

type Tx interface {
	// Gets the Msg.
	GetMsgs() []utils.Msg
}
type SendData struct {
	Tx   StdTx  `json:"tx"`
	Mode string `json:"mode"`
}

// StdTx def
type StdTx struct {
	Msgs       []utils.Msg    `json:"msg"`
	Fee        StdFee         `json:"fee" yaml:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

func (tx StdTx) Route() string { return "hbc" }

// NewStdTx to instantiate an instance
func NewStdTx(msgs []utils.Msg, sigs []StdSignature, memo string, fee StdFee) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

// GetMsgs def
func (tx StdTx) GetMsgs() []utils.Msg { return tx.Msgs }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() error {
	stdSigs := tx.GetSignatures()

	if tx.Fee.Gas > MaxGasWanted {
		return fmt.Errorf(
			"invalid gas supplied; %d > %d", tx.Fee.Gas, MaxGasWanted,
		)
	}

	if len(stdSigs) == 0 {
		return fmt.Errorf("ErrNoSignatures")
	}
	if len(stdSigs) != len(tx.GetSigners()) {
		return fmt.Errorf(
			"wrong number of signers; expected %d, got %d", tx.GetSigners(), len(stdSigs),
		)
	}

	return nil
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigners() []utils.CUAddress {
	var signers []utils.CUAddress
	seen := map[string]bool{}

	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	return signers
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetSignatures returns the signature of signers who signed the Msg.
// CONTRACT: Length returned is same as length of
// pubkeys returned from MsgKeySigners, and the order
// matches.
// CONTRACT: If the signature is missing (ie the Msg is
// invalid), then the corresponding signature is
// .Empty().
func (tx StdTx) GetSignatures() [][]byte {
	sigs := make([][]byte, len(tx.Signatures))
	for i, stdSig := range tx.Signatures {
		sigs[i] = stdSig.Signature
	}
	return sigs
}
