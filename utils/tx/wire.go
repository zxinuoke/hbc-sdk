package tx

import (
	"github.com/bluehelix-chain/hbc-sdk/utils"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

// cdc global variable
var Cdc = amino.NewCodec()

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterConcrete(StdTx{}, "hbtcchain/StdTx", nil)
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	utils.RegisterCodec(cdc)
}

func init() {
	cryptoAmino.RegisterAmino(Cdc)
	RegisterCodec(Cdc)
}
