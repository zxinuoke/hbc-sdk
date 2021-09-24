package hbc

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zxinuoke/hbc-sdk/utils/base58"
)

func TestHbcAddress(t *testing.T) {
	derivedPriStr := "da0bbe0acb8aae423de68a1a59379512d9d92ed453743592d6c3b2bc04252640"
	//pubkeyStr := "eb5ae9872102686ec82d4c6612d030929c3fac296d2bcddca43b2b420810aaaf5915ccec1ad1"

	prikeyByte, _ := hex.DecodeString(derivedPriStr)

	priv := SecpPrivKeyGen(prikeyByte)
	pub := priv.PubKey()
	pubStr := string(hex.EncodeToString(pub.Bytes()))
	fmt.Printf("pubStr %v", pubStr)

	address, _, err := CreateAddress(prikeyByte)
	fmt.Printf("address %v", address)

	valid, _ := CheckAddrValid("HBCb1bg1Y2qxRhVQBUxHE7nWcuKzbM7scrwU")
	fmt.Printf("valid %v", valid)

	msgdata := []byte("my first message")
	signData, err := priv.Sign(msgdata)
	fmt.Printf("signData %v, err %v", signData, err)
}

func TestHbcSign(t *testing.T) {
	derivedPriStr := "80ebd557b994081d9ee8e9896c170a7c4045dcd6aa7b4a231065b4e753707f7d"
	//pubkeyStr := "eb5ae9872102686ec82d4c6612d030929c3fac296d2bcddca43b2b420810aaaf5915ccec1ad1"

	prikeyByte, _ := hex.DecodeString(derivedPriStr)

	priv := SecpPrivKeyGen(prikeyByte)
	pub := priv.PubKey()
	pubStr := string(hex.EncodeToString(pub.Bytes()))
	fmt.Printf("pubStr %v", pubStr)

	msgdata := []byte("my first message")
	signData, err := priv.Sign(msgdata)
	fmt.Printf("signData %v, err %v", signData, err)
}

func TestHbcMultiSign(t *testing.T) {
	//derivedPriStr := "80ebd557b994081d9ee8e9896c170a7c4045dcd6aa7b4a231065b4e753707f7d"
	//derivedPriStr2 := "7ebd557b994081d9ee8e9896c170a7c4045dcd6aa7b4a231065b4e753707f7d"

	derivedPriStr := "01ee5aa673f63fc906fb2dbc191438217c5e3f2646381b5e261be2d1f8479086"
	derivedPriStr2 := "1f118af86fcab84f1a7d5204e0cdda351dcf759498e3550a858b56fc719f5521"
	var derivedPriStrs = []string{derivedPriStr, derivedPriStr2}
	//pubkeyStr := "eb5ae9872102686ec82d4c6612d030929c3fac296d2bcddca43b2b420810aaaf5915ccec1ad1"

	var prikeys [][]byte

	for _, value := range derivedPriStrs {
		prikeyByte, _ := hex.DecodeString(value)
		prikeys = append(prikeys, prikeyByte)
	}

	mulAddress, _, _ := GetMultiAddress(prikeys)

	fmt.Printf("mulAddress %v", mulAddress)
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}

type keyData struct {
	priv string
	pub  string
	addr string
}

var secpDataTable = []keyData{
	{
		priv: "a96e62ed3955e65be32703f12d87b6b5cf26039ecfa948dc5107a495418e5330",
		pub:  "02950e1cdfcb133d6024109fd489f734eeb4502418e538c28481f22bce276f248c",
		addr: "1CKZ9Nx4zgds8tU7nJHotKSDr4a9bYJCa3",
	},
}

func TestPubKeySecp256k1Address(t *testing.T) {
	for _, d := range secpDataTable {
		privB, _ := hex.DecodeString(d.priv)
		pubB, _ := hex.DecodeString(d.pub)
		addrBbz, _, _ := base58.CheckDecode(d.addr)
		addrB := crypto.Address(addrBbz)

		var priv secp256k1.PrivKeySecp256k1
		copy(priv[:], privB)

		pubKey := priv.PubKey()
		pubT, _ := pubKey.(secp256k1.PubKeySecp256k1)
		pub := pubT[:]
		addr := pubKey.Address()

		assert.Equal(t, pub, pubB, "Expected pub keys to match")
		assert.Equal(t, addr, addrB, "Expected addresses to match")
	}
}
