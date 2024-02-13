package utils

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/bitcoinschema/go-bitcoin"
)

// Keys for testing and core signing using Genesis keys
const (
	Key1 = "e785e99372b5ecc65828d120ad7d7c4f8ce7b3496833ad0143cff7eff90dc822"
	Key2 = "7445f20674fe5e051ac5b6712ad250760686cd9a38676d3f92e165380708903e"
	Key3 = "4e7d00a25e93a5736f294a2867b3ec73a8e7bdd88ea35e390c203d52684e5292"
	Key4 = "0f5491fd840b47293ee85a4068b05249ba66958185ca1626442c472ef16cc2aa"
	Key5 = "161a9d680f1952fef2f2572bfae30e7eda047908a180cbb00a16a4ee270abf62"
)

// Mainnet Keys for the Genesis keys
const (
	MainKey1 = "02a1589f2c8e1a4e7cbf28d4d6b676aa2f30811277883211027950e82a83eb2768"
	MainKey2 = "03aec1d40f02ac7f6df701ef8f629515812f1bcd949b6aa6c7a8dd778b748b2433"
	MainKey3 = "03ddb2806f3cc48aa36bd4aea6b9f1c7ed3ffc8b9302b198ca963f15beff123678"
	MainKey4 = "036846e3e8f4f944af644b6a6c6243889dd90d7b6c3593abb9ccf2acb8c9e606e2"
	MainKey5 = "03e45c9dd2b34829c1d27c8b5d16917dd0dc2c88fa0d7bad7bffb9b542229a9304"
)

// SignWithGenesis will sign the data with the genesis keys
func SignWithGenesis(data []byte) ([][]byte, error) {
	return SignWithKeys(data, []string{Key1, Key2, Key3})
}

// SignWithKeys will sign the data with the keys provided
func SignWithKeys(data []byte, keys []string) ([][]byte, error) {
	// var sigs [][]byte
	sigs := make([][]byte, 0)
	var b []byte
	for _, key := range keys {
		s, err := bitcoin.SignMessage(key, hex.EncodeToString(data), true)
		if err != nil {
			return nil, err
		}
		if b, err = base64.StdEncoding.DecodeString(s); err != nil {
			return nil, err
		}
		sigs = append(sigs, b)
	}
	return sigs, nil
}
