package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func doRSAEncrypt(encstring, pubKeyE, pubKeyN string) (string, error) {
	if encstring == "" {
		return "", nil
	}

	block, _ := pem.Decode([]byte(pubKeyN))
	if block == nil {
		return "", errors.New("invalid PEM data")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	encData := base64Encode(encstring)
	chunkSize := 245
	if true { // simulate 'rsapadingtype == "1"'
		chunkSize = 214
	}
	var result []byte
	for i := 0; i < len(encData); i += chunkSize {
		end := i + chunkSize
		if end > len(encData) {
			end = len(encData)
		}
		cipherChunk, err := rsa.EncryptOAEP(
			nil,
			rand.Reader,
			rsaPub,
			[]byte(encData[i:end]),
			nil,
		)
		if err != nil {
			return "", err
		}
		result = append(result, cipherChunk...)
	}
	return hex.EncodeToString(result), nil
}

func dataDecrypt(encrypted, keyStr, ivStr string) (string, error) {
	key, err := hex.DecodeString(keyStr)
	if err != nil {
		return "", err
	}
	iv, err := hex.DecodeString(ivStr)
	if err != nil {
		return "", err
	}
	cipherData, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherData, cipherData)
	return string(cipherData), nil
}
