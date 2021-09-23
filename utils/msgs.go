package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgSend struct {
	FromAddress CUAddress `json:"from_address,omitempty"`
	ToAddress   CUAddress `json:"to_address,omitempty"`
	Amount      Coins     `json:"amount"`
}

var _ Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr CUAddress, amount Coins) MsgSend {
	return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return "bank" }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() error {
	if len(msg.FromAddress) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if len(msg.ToAddress) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Amount inValid")
	}
	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Amount inValid")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(MsgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []CUAddress {
	return []CUAddress{msg.FromAddress}
}

func (msg MsgSend) GetInvolvedAddresses() []CUAddress {
	var addrs []CUAddress
	addrs = append(addrs, msg.FromAddress)
	addrs = append(addrs, msg.ToAddress)
	return addrs
}

type MsgMultiSend struct {
	Inputs  []Input  `protobuf:"bytes,1,rep,name=inputs,proto3" json:"inputs"`
	Outputs []Output `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs"`
}
type Input struct {
	Address CUAddress `json:"address,omitempty"`
	Coins   Coins     `json:"coins"`
}
type Output struct {
	Address CUAddress `json:"address,omitempty"`
	Coins   Coins     ` json:"coins"`
}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

// Route Implements Msg
func (msg MsgMultiSend) Route() string { return "bank" }

// Type Implements Msg
func (msg MsgMultiSend) Type() string { return "multisend" }

// ValidateBasic Implements Msg.
func (msg MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return fmt.Errorf("ErrNoInputs")
	}
	if len(msg.Outputs) == 0 {
		return fmt.Errorf("ErrNoOutputs")
	}

	return ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(MsgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []CUAddress {
	addrs := make([]CUAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

func (msg MsgMultiSend) GetInvolvedAddresses() []CUAddress {
	numOfInputs := len(msg.Inputs)
	numOfOutputs := len(msg.Outputs)
	addrs := make([]CUAddress, numOfInputs+numOfOutputs, numOfInputs+numOfOutputs)
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	for i, out := range msg.Outputs {
		addrs[i+numOfInputs] = out.Address
	}
	return addrs
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if len(in.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "input address missing")
	}
	if !in.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "in.Coins inValid")
	}
	if !in.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "in.Coins inValid")
	}
	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
func NewInput(addr CUAddress, coins Coins) Input {
	return Input{
		Address: addr,
		Coins:   coins,
	}
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if len(out.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "output address missing")
	}
	if !out.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "out.Coins inValid")
	}
	if !out.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "out.Coins inValid")
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
func NewOutput(addr CUAddress, coins Coins) Output {
	return Output{
		Address: addr,
		Coins:   coins,
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) error {
	var totalIn, totalOut Coins

	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}

		totalIn = append(totalIn, in.Coins...)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}

		totalOut = append(totalOut, out.Coins...)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return fmt.Errorf("ErrInputOutputMismatch")
	}

	return nil
}
