package hbc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/zxinuoke/hbc-sdk/utils"
	"github.com/zxinuoke/hbc-sdk/utils/tx"
)

type Hbc struct {
	RestUrl string
}

var (
	DefaultDecimals = 18
	DefaultGasLimit = 2000000
	DefaultFee      = "1000000000000"
	DefaultChainID  = "hbtc-testnet"
	DefaultTokenId  = "hbc"
)

func NewHbc(url, restPort string) (*Hbc, error) {
	if url == "" {
		return nil, errors.New("err NewHbc params")
	}

	restUrl := fmt.Sprintf("%s:%s", url, restPort)

	return &Hbc{
		RestUrl: restUrl,
	}, nil
}

func NewHbcClient(url string) (*Hbc, error) {
	if url == "" {
		return nil, errors.New("err NewBnb params")
	}

	return &Hbc{
		RestUrl: url,
	}, nil
}

func (hbc *Hbc) GetCurrentHeight() (int64, error) {
	var response BlockData
	err := hbc.RequestHbcData("GET", "/blocks/latest", map[string]interface {
	}{}, &response)
	if err != nil {
		return 0, err
	}
	return response.Block.Header.Height.Int64()
}

func (hbc *Hbc) GetBlockData(height int64) (*BlockData, error) {
	var response BlockData
	err := hbc.RequestHbcData("GET", "/blocks/"+strconv.FormatInt(height, 10), map[string]interface {
	}{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (hbc *Hbc) GetTransactionData(hash string) (*TxData, error) {
	var response TxData
	err := hbc.RequestHbcData("GET", "/txs/"+hash, map[string]interface {
	}{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (hbc *Hbc) GetAddressInfo(address string) (*AddressData, error) {
	var response AddressData
	err := hbc.RequestHbcData("GET", "/cu/cus/"+address, map[string]interface {
	}{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (hbc *Hbc) GetCoinBalance(address, coin string) (string, error) {
	if coin == "" {
		coin = "hbc"
	}

	var response AssetData
	err := hbc.RequestHbcData("GET", "/transfer/balances/"+address, map[string]interface {
	}{}, &response)
	if err != nil {
		return "", err
	}

	for _, value := range response.Result.Available {
		if value.Denom == coin {
			return value.Amount, nil
		}
	}

	return "0", errors.New("can not find balance")
}

func createUnsignData(tokenId, fromAddress, toAddress, memo string, amount, fee string, sequence int64) (*tx.StdSignMsg, error) {
	addr1, err := utils.CUAddressFromBase58(fromAddress)
	if err != nil {
		return nil, err
	}
	addr2, err := utils.CUAddressFromBase58(toAddress)
	if err != nil {
		return nil, err
	}
	amountBigInt, ok := sdk.NewIntFromString(amount)
	if !ok {
		return nil, errors.New("error send amount")
	}
	coins := utils.NewCoins(utils.NewCoin(tokenId, amountBigInt))
	msg := utils.NewMsgSend(addr1, addr2, coins)

	feeBigInt, ok := sdk.NewIntFromString(fee)
	if !ok {
		return nil, errors.New("error send fee")
	}

	defaultFeeBigInt, ok := sdk.NewIntFromString(DefaultFee)
	if !ok {
		return nil, errors.New("error  defaultFeeBigInt")
	}
	if feeBigInt.LT(defaultFeeBigInt) {
		feeBigInt = defaultFeeBigInt
	}

	feecoins := sdk.NewCoins(sdk.NewCoin(DefaultTokenId, feeBigInt))
	feeData := tx.NewStdFee(uint64(DefaultGasLimit), feecoins)

	var msgs []utils.Msg
	msgs = append(msgs, msg)

	signMsg := tx.StdSignMsg{
		ChainID:  DefaultChainID,
		Sequence: sequence,
		Memo:     memo,
		Msgs:     msgs,
		Fee:      feeData,
	}
	for _, m := range signMsg.Msgs {
		if err := m.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	return &signMsg, nil
}

func CreateTransaction(tokenId string, fromPriKey []byte, fromAddress, toAddress, memo string, amount, fee string, sequence int64) ([]byte, error) {
	signMsg, err := createUnsignData(tokenId, fromAddress, toAddress, memo, amount, fee, sequence)
	if err != nil {
		return nil, err
	}
	priv := SecpPrivKeyGen(fromPriKey)

	signData, err := priv.Sign(signMsg.Bytes())
	if err != nil {
		return nil, err
	}

	pubT, ok := priv.PubKey().(secp256k1.PubKeySecp256k1)
	if !ok {
		return nil, err
	}

	sig := tx.StdSignature{
		PubKey:    pubT,
		Signature: signData,
	}

	stdTx := tx.NewStdTx(signMsg.Msgs, []tx.StdSignature{sig}, signMsg.Memo, signMsg.Fee)

	err = stdTx.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sendData := tx.SendData{Tx: stdTx, Mode: "sync"}

	bz, err := tx.Cdc.MarshalJSON(&sendData)
	if err != nil {
		return nil, err
	}

	return bz, err
}

func CreateMultiTransaction(tokenId string, fromPriKey []byte, pubkeys []crypto.PubKey, fromAddress, toAddress, memo string, amount, fee string, sequence int64) ([]byte, error) {
	signMsg, err := createUnsignData(tokenId, fromAddress, toAddress, memo, amount, fee, sequence)
	if err != nil {
		return nil, err
	}
	priv := SecpPrivKeyGen(fromPriKey)

	signData, err := priv.Sign(signMsg.Bytes())
	if err != nil {
		return nil, err
	}

	pubT, ok := priv.PubKey().(secp256k1.PubKeySecp256k1)
	if !ok {
		return nil, err
	}

	mpk := multisig.PubKeyMultisigThreshold{
		K:       2,
		PubKeys: pubkeys,
	}

	multisigSig := multisig.NewMultisig(len(pubkeys))
	if err := multisigSig.AddSignatureFromPubKey(signData, pubT, mpk.PubKeys); err != nil {
		return nil, err
	}

	encodeData, err := json.Marshal(multisigSig)
	if err != nil {
		return nil, err
	}

	return encodeData, nil
}

func MergeMultiSign(tokenId string, txData []byte, fromPriKey []byte, pubkeys []crypto.PubKey, fromAddress, toAddress, memo string, amount, fee string, sequence int64) ([]byte, error) {
	var multisigSig multisig.Multisignature

	err := json.Unmarshal(txData, &multisigSig)
	if err != nil {
		return nil, err
	}

	signMsg, err := createUnsignData(tokenId, fromAddress, toAddress, memo, amount, fee, sequence)
	if err != nil {
		return nil, err
	}
	priv := SecpPrivKeyGen(fromPriKey)

	signData, err := priv.Sign(signMsg.Bytes())
	if err != nil {
		return nil, err
	}

	pubT, ok := priv.PubKey().(secp256k1.PubKeySecp256k1)
	if !ok {
		return nil, err
	}
	mpk := multisig.PubKeyMultisigThreshold{
		K:       2,
		PubKeys: pubkeys,
	}

	if err := multisigSig.AddSignatureFromPubKey(signData, pubT, mpk.PubKeys); err != nil {
		return nil, err
	}

	ok = mpk.VerifyBytes(signMsg.Bytes(), multisigSig.Marshal())
	if !ok {
		return nil, errors.New("verify sign failed")
	}

	sig := tx.StdSignature{PubKey: mpk, Signature: multisigSig.Marshal()}
	stdTx := tx.NewStdTx(signMsg.Msgs, []tx.StdSignature{sig}, signMsg.Memo, signMsg.Fee)
	err = stdTx.ValidateBasic()
	if err != nil {
		return nil, err
	}

	sendData := tx.SendData{Tx: stdTx, Mode: "sync"}

	bz, err := tx.Cdc.MarshalJSON(&sendData)
	if err != nil {
		return nil, err
	}

	return bz, err
}

func (hbc *Hbc) SendSignedTx(txData []byte) (string, error) {
	var response TxResponse

	err := hbc.PostHbcData("/txs", txData, &response)
	if err != nil {
		return "", err
	}

	return response.TxHash, nil
}

func (hbc *Hbc) PostHbcData(requestPath string, postData []byte, model interface{}) error {
	if len(postData) == 0 {
		return errors.New("no postData")
	}
	client := http.Client{}

	requestString := string(postData)
	req, err := http.NewRequest("POST", hbc.RestUrl+requestPath, strings.NewReader(requestString))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	stringBody := string(bytes)

	errResp := &BaseResponse{}
	err = json.Unmarshal(bytes, errResp)
	if nil == err && errResp.Error != "" {
		return fmt.Errorf("http request err:%v", stringBody)
	}

	err = json.Unmarshal(bytes, &model)
	if err != nil {
		return err
	}

	return nil
}

func (hbc *Hbc) RequestHbcData(method, requestPath string, args map[string]interface{}, model interface{}) error {
	client := http.Client{}
	var req *http.Request

	var requestString string
	var err error
	var mjson []byte

	if method == "POST" {
		mjson, err = json.Marshal(args)
		if err != nil {
			return err
		}
		requestString = string(mjson)

		req, err = http.NewRequest("POST", hbc.RestUrl+requestPath, strings.NewReader(requestString))
	} else if method == "GET" {
		requestString = requestPath
		if len(args) > 0 {
			requestString = requestString + "?"
			for key := range args { //取map中的值
				requestString = requestString + key + "=" + args[key].(string) + "&"
			}
			requestString = strings.TrimSuffix(requestString, "&")
		}

		req, err = http.NewRequest("GET", hbc.RestUrl+requestString, nil)
	}
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	errResp := &BaseResponse{}
	err = json.Unmarshal(bytes, errResp)
	if nil == err && errResp.Error != "" {
		return fmt.Errorf("http request err:%v", string(bytes))
	}

	err = json.Unmarshal(bytes, &model)
	if err != nil {
		return err
	}

	return nil
}
