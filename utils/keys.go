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
