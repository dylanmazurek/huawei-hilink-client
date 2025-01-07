package crypto

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"math"
	"math/big"
)

func RSAEncrypt(inputStr string, pubKeyN string) (*string, error) {
	N, ok := new(big.Int).SetString(pubKeyN, 16)
	if !ok {
		return nil, nil
	}

	pubKey := rsa.PublicKey{N: N, E: 65537}

	randomStr := []byte{213, 248, 224, 85, 93, 22, 147, 239, 235, 92, 51, 32, 155, 142, 104, 70, 210, 82, 86, 165}
	base64InputStr := base64.StdEncoding.EncodeToString([]byte(inputStr))

	sectionCount := int(math.Ceil(float64(len(base64InputStr)) / 214))
	if sectionCount == 0 {
		sectionCount = 1
	}

	var encodedStr []byte
	for i := 0; i < sectionCount; i++ {
		start := i * 214
		end := start + 214

		if end > len(base64InputStr) {
			end = len(base64InputStr)
		}

		encdata := base64InputStr[start:end]
		encrypted, err := rsa.EncryptOAEP(sha1.New(), bytes.NewReader(randomStr), &pubKey, []byte(encdata), nil)

		if err != nil {
			return nil, err
		}

		encodedStr = append(encodedStr, encrypted...)
	}

	outStr := hex.EncodeToString(encodedStr)
	return &outStr, nil
}
