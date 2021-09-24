package hbc

import (
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zxinuoke/hbc-sdk/utils"
)

func CheckHbcPrivkeyAddr(privKey []byte, addrInput string) (bool, error) {
	address, _, err := CreateAddress(privKey)
	if err != nil {
		return false, err
	}

	if strings.Compare(address, addrInput) != 0 {
		return false, errors.New("address not matched with privkey")
	}
	return true, nil
}

func SecpPrivKeyGen(bz []byte) crypto.PrivKey {
	var bzArr [32]byte
	copy(bzArr[:], bz)
	return secp256k1.PrivKeySecp256k1(bzArr)
}

func CreateAddress(bz []byte) (string, []byte, error) {
	priv := SecpPrivKeyGen(bz)
	pub := priv.PubKey()

	addr := pub.Address()

	acc := utils.CUAddress(addr)

	priT, ok := priv.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return "", nil, errors.New("parse address error")
	}

	return acc.String(), priT[:], nil
}

func GetMultiAddress(privateKeys [][]byte) (string, []byte, error) {
	var pks []crypto.PubKey

	for _, pk := range privateKeys {
		priv := SecpPrivKeyGen(pk)
		pub := priv.PubKey()
		pks = append(pks, pub)
	}

	mpk := multisig.PubKeyMultisigThreshold{K: 2, PubKeys: pks}

	addr := mpk.Address()

	acc := utils.CUAddress(addr)

	return acc.String(), mpk.Bytes(), nil
}

func GetMultiPubs(pks []byte) ([]crypto.PubKey, error) {
	cdc := codec.New()
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{}, "tendermint/PubKeyMultisigThreshold", nil)
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)

	var mpks multisig.PubKeyMultisigThreshold
	err := cdc.UnmarshalBinaryBare(pks, &mpks)
	if err != nil {
		return nil, err
	}
	return mpks.PubKeys, nil
}

func CheckAddrValid(addr string) (bool, error) {
	_, err := utils.CUAddressFromBase58(addr)
	if err != nil {
		return false, err
	}
	return true, nil
}
