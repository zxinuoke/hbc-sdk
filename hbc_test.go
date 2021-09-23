package hbc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/types"
)

const (
	testUrl  = "http://xxxxx"
	testPort = "1317"
)

func TestHbcCommon(t *testing.T) {
	/*	data := "xgLwYl3uCkgqLIf6CiAKFFHQDFP1ovY3g1F+vk68fEjBIsEJEggKA0JOQhCQThIgChSOpw19LqihS6KzPRjV371vrgpuqBIICgNCTkIQkE4S6gEKViLB9+IIAhIm61rphyEDki0s+BmCtRuy4wc6DeTUneeg3C2t7jU0tMV/p2lo/mgSJuta6YchA2c+8e4xcKsZhH/CDMPty+QPdgE+VKL/eYAY2wWJA+ZfEosBCgUIAhIBwBJAfpNlO4LfEEwVjv7JyCLbkykgRLIVkVFWBf1OOwQCirYUxkjLbxKwFB9K91ukn5bsnqVh0gFidVDvDAPTz9CiKhJA269YzyXpPI+lyqPF/0E5QpFjEnzxLhNkVb7U0JIFcHYEdNe9hK++vqW531V8s+SeHxa8CUxQ/iXH8Mb2H4924RiMwx4aCTEwNDQxNjQ1MQ=="
		byteData, err := base64.StdEncoding.DecodeString(data)
		tx := TxData(byteData)
		fmt.Printf("hash %x, string %v , string byteData %v \n", tx.Hash(), tx.String(), string(byteData))
		fmt.Printf("err %v \n", err)*/

	pri, _ := hex.DecodeString("pri_hex_1")
	address, _, _ := CreateAddress(pri)
	fmt.Printf("address %v \n", address)

	client, _ := NewHbc(testUrl, testPort)

	gas, err := GetHbcGas()
	fmt.Printf("gas %v \n", gas)
	fmt.Printf("err %v \n", err)

	info, err := client.GetCurrentHeight()
	fmt.Printf("GetCurrentHeight %v \n", info)
	fmt.Printf("err %v \n", err)

	blockData, err := client.GetBlockData(646737)
	fmt.Printf("blockData %v \n", blockData)

	fmt.Printf("err %v \n", err)
	txData, err := client.GetTransactionData("B3B3498D178A5EB86EFED7FC287754429754263BAE39032785AE7625270E539B")
	fmt.Printf("txData %v \n", txData)
	for _, value := range blockData.Block.Data.Txs {
		byteData, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			continue
		}
		tx := types.Tx(byteData)
		txData, err := client.GetTransactionData(hex.EncodeToString(tx.Hash()))
		fmt.Printf("txData %v \n", txData)
		fmt.Printf("err %v \n", err)
	}

	txData, err = client.GetTransactionData("9F7E147411FCFEDD9C44836C74BE061D8452A2708E01F6B6587C013586EF1FE1")
	fmt.Printf("txData %v \n", txData)
	fmt.Printf("err %v \n", err)

	addressData, err := client.GetAddressInfo("HBCb1bg1Y2qxRhVQBUxHE7nWcuKzbM7scrwU")
	fmt.Printf("addressData %v \n", addressData)
	fmt.Printf("err %v \n", err)

	balance, err := client.GetCoinBalance("f1kgyowkdn3pz5u6fecb5buyby6cf2s7wjwvhdcva", "KIWI")
	fmt.Printf("balance %v \n", balance)
	fmt.Printf("err %v \n", err)
}

func TestHbc_CreateAndSignTransaction(t *testing.T) {
	client, _ := NewHbc(testUrl, testPort)

	derivedPriStr := "pri_hex_1"
	prikeyByte, _ := hex.DecodeString(derivedPriStr)
	address, _, err := CreateAddress(prikeyByte)
	if err != nil {
		return
	}
	fmt.Printf("address %v  \n", address)

	fromAddressInfo, err := client.GetAddressInfo("HBCb1bg1Y2qxRhVQBUxHE7nWcuKzbM7scrwU")
	if err != nil {
		return
	}

	coinAsset, err := client.GetCoinBalance("HBCb1bg1Y2qxRhVQBUxHE7nWcuKzbM7scrwU", "HBCGLw2dz8hXHRaotJEA9QHxdRYybqKV7UuG")
	if err != nil {
		return
	}
	fmt.Printf("coinAsset %v  \n", coinAsset)

	sequence, err := fromAddressInfo.Result.Value.Sequence.Int64()
	if err != nil {
		return
	}

	txData, err := CreateTransaction("HBCGLw2dz8hXHRaotJEA9QHxdRYybqKV7UuG", prikeyByte, "HBCb1bg1Y2qxRhVQBUxHE7nWcuKzbM7scrwU", "HBCTeUXgzx8eenRXmd6ztAJe4xdQmjMFUV4t", "", "100000000", DefaultFee, sequence)
	if err != nil {
		return
	}

	txDataStr := string(txData)
	fmt.Printf("txDataStr %v  \n", txDataStr)
	hash, err := client.SendSignedTx(txData)
	if err != nil {
		return
	}
	fmt.Printf("hash %v  \n", hash)
}

func TestHbc_multiSign(t *testing.T) {
	client, _ := NewHbc(testUrl, testPort)
	derivedPriStr := "pri_hex_1"
	derivedPriStr2 := "pri_hex_2"

	var derivedPriStrs = []string{derivedPriStr, derivedPriStr2}

	var prikeys [][]byte

	for _, value := range derivedPriStrs {
		prikeyByte, _ := hex.DecodeString(value)
		prikeys = append(prikeys, prikeyByte)
	}

	mulAddress, pks, _ := GetMultiAddress(prikeys) //HBCTeUXgzx8eenRXmd6ztAJe4xdQmjMFUV4t

	cdc := codec.New()
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{}, "tendermint/PubKeyMultisigThreshold", nil)
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)

	var mpks multisig.PubKeyMultisigThreshold
	err := cdc.UnmarshalBinaryBare(pks, &mpks)
	if err != nil {
		return
	}

	fmt.Printf("mulAddress %v", mulAddress)

	//address1, _, err := CreateAddress(prikeys[0])
	fromAddressInfo, err := client.GetAddressInfo(mulAddress)
	sequence, err := fromAddressInfo.Result.Value.Sequence.Int64()
	//address2, _, err := CreateAddress(prikeys[1])
	fromAddressInfo2, err := client.GetAddressInfo(mulAddress)
	sequence2, err := fromAddressInfo2.Result.Value.Sequence.Int64()
	txData, err := CreateMultiTransaction("HBCGLw2dz8hXHRaotJEA9QHxdRYybqKV7UuG", prikeys[0], mpks.PubKeys, "HBCTeUXgzx8eenRXmd6ztAJe4xdQmjMFUV4t", "HBCgKep1AQKT1x9KhsDUThyRzkRMkYYoCGT8", "104416451", "10000000", DefaultFee, sequence)

	txDataMerge, err := MergeMultiSign("HBCGLw2dz8hXHRaotJEA9QHxdRYybqKV7UuG", txData, prikeys[1], mpks.PubKeys, "HBCTeUXgzx8eenRXmd6ztAJe4xdQmjMFUV4t", "HBCgKep1AQKT1x9KhsDUThyRzkRMkYYoCGT8", "104416451", "10000000", DefaultFee, sequence2)

	txDataStr := string(txDataMerge)
	fmt.Printf("txDataStr %v  \n", txDataStr)
	hash, err := client.SendSignedTx(txDataMerge)
	if err != nil {
		return
	}
	fmt.Printf("hash %v  \n", hash)
}
