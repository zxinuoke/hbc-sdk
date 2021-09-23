package utils

import (
	"github.com/tendermint/go-amino"
)

var MsgCdc = amino.NewCodec()

func RegisterCodec(cdc *amino.Codec) {
	//Must use cosmos-sdk.
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterConcrete(MsgSend{}, "hbtcchain/transfer/MsgSend", nil)
}

func init() {
	RegisterCodec(MsgCdc)
}
