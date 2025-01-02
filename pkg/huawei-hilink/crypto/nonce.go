package crypto

import "crypto/rand"

func newNonce(keySizeBytes int) ([]byte, error) {
	if keySizeBytes <= 0 {
		return nil, ErrInvalidKeySize
	}

	bytes := make([]byte, keySizeBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
