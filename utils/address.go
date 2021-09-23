package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bluehelix-chain/hbc-sdk/utils/base58"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"gopkg.in/yaml.v2"
)

const (
	// Constants defined here are the defaults value for address.
	// You can use the specific values for your project.
	// Add the follow lines to the `main()` of your server.
	//
	//	config := sdk.GetConfig()
	//	config.SetBech32PrefixForAccount(yourBech32PrefixAccAddr, yourBech32PrefixAccPub)
	//	config.SetBech32PrefixForValidator(yourBech32PrefixValAddr, yourBech32PrefixValPub)
	//	config.SetBech32PrefixForConsensusNode(yourBech32PrefixConsAddr, yourBech32PrefixConsPub)
	//	config.SetCoinType(yourCoinType)
	//	config.SetFullFundraiserPath(yourFullFundraiserPath)
	//	config.Seal()

	// AddrLen defines a valid address length
	AddrLen = 20
	// Bech32PrefixAccAddr defines the Bech32 prefix of an CU's address
	Bech32MainPrefix = "hbc"

	// bht in https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	CoinType = 496

	// BIP44Prefix is the parts of the BIP44 HD path that are fixed by
	// what we used during the fundraiser.
	FullFundraiserPath = "44'/496'/0'/0/0"

	// PrefixAccount is the prefix for CU keys
	PrefixAccount = "acc"
	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// PrefixAddress is the prefix for addresses
	PrefixAddress = "addr"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an CU's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an CU's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)

type CUType int

const (
	CUTypeUser CUType = 0x1 //用户地址
	CUTypeOp   CUType = 0x2 //运营地址
	CUTypeORG  CUType = 0x3 //机构地址
)

// Address is a common interface for different types of addresses used by the SDK
type Address interface {
	Equals(Address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// Ensure that different address types implement the interface
var _ Address = CUAddress{}

var _ yaml.Marshaler = CUAddress{}

// ----------------------------------------------------------------------------
// CU
// ----------------------------------------------------------------------------

// CUAddress a wrapper around bytes meant to represent an CU address.
// When marshaled to a string or JSON, it uses Base58.
// Implement address interface
type CUAddress []byte

// CUAddressFromHex creates an CUAddress from a hex string.
func CUAddressFromHex(address string) (addr CUAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding hex address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return CUAddress(bz), nil
}

// CUAddressFromBase58 creates an AccAddress from a base58 string prefixed with "HBT".
func CUAddressFromBase58(address string) (addr CUAddress, err error) {
	// blank input get CUAddress{} without error
	if len(strings.TrimSpace(address)) == 0 {
		return CUAddress{}, nil
	}

	// TODO uncomment
	if len(strings.TrimSpace(address)) < base58.AddrPrefixLen || address[:base58.AddrPrefixLen] != base58.AddrStrPrefix {
		return nil, errors.New(fmt.Sprintf("invalid cuaddress:%v with prefixed !=HBC", address))
	}

	bz, version, err := base58.CheckDecode(address)

	if err != nil {
		return CUAddress{}, err
	}

	if len(bz) != (AddrLen + base58.AddrPrefixLen - 1) { //?
		return nil, errors.New("Incorrect address length")
	}

	prefix := make([]byte, 0, base58.AddrPrefixLen)
	prefix = append(prefix, version)
	prefix = append(prefix, bz[:base58.AddrPrefixLen-1]...)

	if !bytes.Equal(prefix, base58.AddrBytePrefix) {
		return CUAddress{}, errors.New("string is not prefixed with `HBC`")
	}
	return CUAddress(bz[2:]), nil
}

func CUAddressFromPubKey(pubKey crypto.PubKey) CUAddress {
	return CUAddress(pubKey.Address().Bytes())
}

func CUAddressFromByte(b []byte) CUAddress {
	if len(b) != AddrLen {
		return CUAddress{}
	}
	return CUAddress(b)
}

// Returns boolean for whether CUAddress equal to another address
func (ca CUAddress) Equals(ca2 Address) bool {
	if ca.Empty() && ca2.Empty() {
		return true
	}

	return bytes.Equal(ca.Bytes(), ca2.Bytes())
}

// Returns boolean for whether an CUAddress is empty
func (ca CUAddress) Empty() bool {
	if ca == nil {
		return true
	}

	ca2 := CUAddress{}
	return bytes.Equal(ca.Bytes(), ca2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (ca CUAddress) Marshal() ([]byte, error) {
	return ca, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (ca *CUAddress) Unmarshal(data []byte) error {
	*ca = data
	return nil
}

// MarshalJSON marshals to JSON using Base58.
func (ca CUAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.String())
}

// UnmarshalJSON unmarshals from JSON assuming Base58 encoding.
func (ca *CUAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	ca2, err := CUAddressFromBase58(s)
	if err != nil {
		return err
	}

	*ca = ca2
	return nil
}

// Bytes returns the raw address bytes.
func (ca CUAddress) Bytes() []byte {
	return ca
}

// String implements the Stringer interface.

func (ca CUAddress) String() string {
	if len(ca) != AddrLen {
		return ""
	}
	b := make([]byte, 0, base58.AddrPrefixLen+len(ca)+4)
	b = append(b, base58.AddrBytePrefix...) //add 'HBC' prefix
	b = append(b, ca[:]...)
	cksum := base58.Checksum(b)
	b = append(b, cksum[:]...)
	return base58.Encode(b)
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (ca CUAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(ca.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", ca)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(ca))))
	}
}

func (ca CUAddress) IsValidAddr() bool {
	if len(ca) != AddrLen {
		return false
	}
	return IsValidAddr(ca.String())
}

func NewCUAddress() CUAddress {
	pubKey := secp256k1.GenPrivKey().PubKey()
	return CUAddress(pubKey.Address())
}

func PubkeyToString(pubkey crypto.PubKey) string {
	return "BHPubKey:" + base58.Encode(pubkey.Bytes())
}

func IsValidAddr(addr string) bool {
	if len(addr) <= base58.AddrPrefixLen {
		return false
	}
	if addr[:3] != base58.AddrStrPrefix {
		return false
	}
	_, _, err := base58.CheckDecode(addr)
	return err == nil
}

type CUAddressList []CUAddress

func (l CUAddressList) Len() int           { return len(l) }
func (l CUAddressList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l CUAddressList) Less(i, j int) bool { return bytes.Compare(l[i], l[j]) == -1 }

func (l CUAddressList) Join() string {
	l2 := make([]string, len(l))
	for i, t := range l {
		l2[i] = t.String()
	}
	return strings.Join(l2, ",")
}

// Any is a method on CUAddressList that returns true if at least one member of the list satisfies a function. It returns false if the list is empty.
func (l CUAddressList) Any(f func(CUAddress) bool) bool {
	for _, t := range l {
		if f(t) {
			return true
		}
	}
	return false
}

func (l CUAddressList) Contains(target Address) bool {
	return l.Any(func(address CUAddress) bool {
		return address.Equals(target)
	})
}

// MarshalYAML marshals to YAML using Bech32.
func (ca CUAddress) MarshalYAML() (interface{}, error) {
	return ca.String(), nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (ca *CUAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	ca2, err := CUAddressFromBase58(s)
	if err != nil {
		return err
	}
	*ca = ca2
	return nil
}
