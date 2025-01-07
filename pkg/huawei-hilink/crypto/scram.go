package crypto

import (
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/xdg-go/pbkdf2"
)

const (
	keySizeBytes    = 32
	clientKeyPrefix = "Client Key"
	serverKeyPrefix = "Server Key"
)

type Scram struct {
	hasher       func() hash.Hash
	keySizeBytes int

	firstNonce []byte
	nonce      []byte

	password       []byte
	saltedPassword []byte

	salt       []byte
	scarmSalt  []byte
	iterations int

	clientKey []byte
	serverKey []byte
	storedKey []byte
	signature []byte
}

func NewScram(opts ...Option) (*Scram, error) {
	clientOptions := DefaultOptions()
	for _, opt := range opts {
		opt(&clientOptions)
	}

	scram := &Scram{
		hasher:       clientOptions.hasher,
		keySizeBytes: clientOptions.keySizeBytes,

		password:   []byte(clientOptions.password),
		iterations: clientOptions.iterations,
	}

	if clientOptions.nonce != nil && clientOptions.finalNonce != nil {
		firstNonce, err := hex.DecodeString(*clientOptions.nonce)
		if err != nil {
			return nil, err
		}

		scram.firstNonce = firstNonce

		finalNonce, err := hex.DecodeString(*clientOptions.finalNonce)
		if err != nil {
			return nil, err
		}

		scram.nonce = finalNonce
	} else {
		newNonce, err := newNonce(scram.keySizeBytes)
		if err != nil {
			return nil, err
		}

		scram.nonce = newNonce
	}

	return scram, nil
}

func (s *Scram) SetSalt(salt string) error {
	if salt == "" || len(salt) == 0 {
		return fmt.Errorf("salt cannot be empty")
	}

	if s.password == nil || len(s.password) == 0 {
		return fmt.Errorf("password cannot be nil")
	}

	scarmSalt, err := hex.DecodeString(salt)
	if err != nil {
		return err
	}

	s.salt = []byte(salt)
	s.scarmSalt = scarmSalt

	dk := pbkdf2.Key(s.password, scarmSalt, s.iterations, s.keySizeBytes, s.hasher)
	s.saltedPassword = dk

	return nil
}

func (s *Scram) GetNonce() ([]byte, error) {
	if s.nonce == nil {
		return nil, fmt.Errorf("nonce cannot be nil")
	}

	return s.nonce, nil
}

func (s *Scram) SetNonce(nonce string) error {
	if nonce == "" {
		return fmt.Errorf("nonce cannot be empty")
	}

	s.firstNonce = s.nonce

	s.nonce = []byte(nonce)
	return nil
}

func (s *Scram) SetPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	s.password = []byte(password)

	return nil
}

func (s *Scram) String() string {
	var scramString string
	scramString += fmt.Sprintf("scramSalt: %s\n", hex.EncodeToString(s.scarmSalt))
	scramString += fmt.Sprintf("saltedPassword: %s\n", hex.EncodeToString(s.saltedPassword))

	scramString += fmt.Sprintf("clientKey: %s\n", hex.EncodeToString(s.clientKey))
	scramString += fmt.Sprintf("serverKey: %s\n", hex.EncodeToString(s.serverKey))
	scramString += fmt.Sprintf("storedKey: %s\n", hex.EncodeToString(s.storedKey))
	scramString += fmt.Sprintf("signature: %s\n", hex.EncodeToString(s.signature))

	scramString += fmt.Sprintf("nonce: %s\n", s.firstNonce)
	scramString += fmt.Sprintf("finalNonce: %s\n", s.nonce)
	return scramString
}
