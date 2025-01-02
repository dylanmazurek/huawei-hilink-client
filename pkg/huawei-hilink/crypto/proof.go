package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func (s *Scram) initClientKey() error {
	h := hmac.New(s.hasher, []byte(clientKeyPrefix))

	_, err := h.Write(s.saltedPassword)
	if err != nil {
		return err
	}

	s.clientKey = h.Sum(nil)
	return nil
}

func (s *Scram) initServerKey() error {
	h := hmac.New(s.hasher, []byte(serverKeyPrefix))

	_, err := h.Write(s.saltedPassword)
	if err != nil {
		return err
	}

	s.serverKey = h.Sum(nil)
	return nil
}

func (s *Scram) initStoredKey() error {
	hasher := sha256.New()
	_, err := hasher.Write(s.clientKey)
	if err != nil {
		return err
	}

	s.storedKey = hasher.Sum(nil)

	return nil
}

func (s *Scram) initSignature() error {
	if s.firstNonce == nil {
		return fmt.Errorf("first nonce cannot be nil")
	}

	if s.nonce == nil {
		return fmt.Errorf("nonce cannot be nil")
	}

	authMessage := []byte(fmt.Sprintf("%s,%s,%s", hex.EncodeToString(s.firstNonce), s.nonce, s.nonce))
	h := hmac.New(s.hasher, authMessage)

	_, err := h.Write(s.storedKey)
	if err != nil {
		return err
	}

	s.signature = h.Sum(nil)

	return nil
}

func (s *Scram) ClientProof() ([]byte, error) {
	err := s.initClientKey()
	if err != nil {
		return nil, err
	}

	err = s.initServerKey()
	if err != nil {
		return nil, err
	}

	err = s.initStoredKey()
	if err != nil {
		return nil, err
	}

	err = s.initSignature()
	if err != nil {
		return nil, err
	}

	clientKeyLen := len(s.clientKey)
	proof := make([]byte, clientKeyLen)
	for i := 0; i < clientKeyLen; i++ {
		proof[i] = s.clientKey[i] ^ s.signature[i]
	}

	return proof, nil
}
