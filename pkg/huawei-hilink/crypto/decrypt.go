package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

func utf8to16(str string) string {
	// Replicates JS "utf8to16" conversion
	var out []rune
	for i := 0; i < len(str); {
		b := str[i]
		switch b >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			out = append(out, rune(b))
			i++
		case 12, 13:
			if i+1 < len(str) {
				b2 := str[i+1]
				r := ((rune(b) & 0x1F) << 6) | (rune(b2) & 0x3F)
				out = append(out, r)
				i += 2
			} else {
				i++
			}
		case 14:
			if i+2 < len(str) {
				b2 := str[i+1]
				b3 := str[i+2]
				r := ((rune(b) & 0x0F) << 12) | ((rune(b2) & 0x3F) << 6) | (rune(b3) & 0x3F)
				out = append(out, r)
				i += 3
			} else {
				i++
			}
		default:
			i++
		}
	}
	return string(out)
}

func pkcs7Pad(data []byte) []byte {
	padLen := aes.BlockSize - (len(data) % aes.BlockSize)
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func bytesToLatin1String(decrypted []byte) (string, error) {
	decoder := charmap.ISO8859_1.NewDecoder()
	latinStr, err := decoder.String(string(decrypted))
	if err != nil {
		return "", err
	}
	return latinStr, nil
}

func AESDecryptCBC(scram Scram, nonce, salt, data string) (*string, error) {
	scram.SetPassword(nonce)
	scram.SetSalt(salt)

	keyStr := hex.EncodeToString(scram.saltedPassword)

	aesKey := keyStr[0:32]
	aesIV := keyStr[32:48]

	log.Info().Msgf("aes [key: %s] [iv: %s]", aesKey, aesIV)

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(aesIV))

	encryptedHexBytes, err := hex.DecodeString(data)
	log.Info().Msgf("encrypted hex str: %s", encryptedHexBytes)

	base64Data := base64.StdEncoding.EncodeToString(encryptedHexBytes)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("encrypted base64 str: %s", base64Data)

	paddedData := pkcs7Pad([]byte(base64Data))

	outputData := make([]byte, len(paddedData))
	mode.CryptBlocks(outputData, []byte(paddedData))
	decryptedLatin, err := bytesToLatin1String(outputData)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("decrypted latin: %s", decryptedLatin)

	decrypted := utf8to16(decryptedLatin)
	log.Info().Msgf("decrypted utf8: %s", decrypted)

	output := hex.EncodeToString([]byte(decrypted))
	log.Info().Msgf("decrypted: %s", output)

	return &output, nil
}
