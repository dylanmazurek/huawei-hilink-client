package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
)

func TestNewScram(t *testing.T) {
	tests := []struct {
		name       string
		password   string
		nonce      *string
		finalNonce *string

		wantErr error
	}{
		{
			name:     "Basic initialization",
			password: "testpass",
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  errors.New("password cannot be empty"),
		},
		{
			name:       "With nonce and final nonce",
			password:   "testpass",
			nonce:      ptr.String("testnonce"),
			finalNonce: ptr.String("finalnonce"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scramOpts := []Option{
				WithPassword(tt.password),
				WithHasher(sha256.New),
			}

			if tt.nonce != nil {
				scramOpts = append(scramOpts, WithStaticNonce(*tt.nonce, *tt.finalNonce))
			}

			scram, err := NewScram(scramOpts...)
			if err != nil {
				if tt.wantErr == nil {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
				}
			} else {
				assert.NotEmpty(t, scram)
				if tt.nonce != nil {
					assert.Equal(t, *tt.nonce, string(scram.nonce))
				}
				if tt.finalNonce != nil {
					assert.Equal(t, *tt.finalNonce, string(scram.nonce))
				}
			}
		})
	}
}

func TestSetSalt(t *testing.T) {
	tests := []struct {
		name       string
		password   string
		salt       string
		iterations int
		wantErr    error
	}{
		{
			name:       "Valid salt",
			password:   "testpass",
			salt:       "73616c74",
			iterations: 100,
		},
		{
			name:     "Empty salt",
			password: "testpass",
			salt:     "",
			wantErr:  errors.New("salt cannot be empty"),
		},
		{
			name:       "Invalid hex salt",
			password:   "testpass",
			salt:       "invalidhex",
			iterations: 100,
			wantErr:    errors.New("encoding/hex: invalid byte: U+0069 'i'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scramOpts := []Option{
				WithPassword(tt.password),
				WithHasher(sha256.New),
				WithIterations(tt.iterations),
			}

			scram, err := NewScram(scramOpts...)
			assert.NoError(t, err)

			err = scram.SetSalt(tt.salt)
			if err != nil {
				if tt.wantErr == nil {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
				}
			} else {
				assert.NotEmpty(t, scram.password)
				assert.NotEmpty(t, scram.salt)
				assert.NotEmpty(t, scram.scarmSalt)
				assert.NotEmpty(t, scram.saltedPassword)
			}
		})
	}
}

func TestClientProof(t *testing.T) {
	tests := []struct {
		name       string
		nonce      string
		finalNonce string
		salt       string
		password   string
		iterations int

		wantErr         error
		wantClientProof string
	}{
		{
			name:       "Valid client proof",
			nonce:      "868f304eec646840959cb3830c0fe352bbc3b24f264ddd652764c8ad2c5fabe3",
			finalNonce: "868f304eec646840959cb3830c0fe352bbc3b24f264ddd652764c8ad2c5fabe31RYXGq1AJapOUsCb57h0bT5oCTsgDIDu",
			password:   "testpass",
			salt:       "4ec91c5b362986cec5ef6db10eb4bb5b1c54eef1eff89006cf00a8993226fabd",
			iterations: 100,

			wantClientProof: "a79ce371acb04daeca2f4f496a9906df970794885ed764760316c54e2433de46",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scramOpts := []Option{
				WithPassword(tt.password),
				WithStaticNonce(tt.nonce, tt.finalNonce),
				WithHasher(sha256.New),
				WithIterations(tt.iterations),
			}

			scram, err := NewScram(scramOpts...)
			assert.NoError(t, err)

			err = scram.SetSalt(tt.salt)
			assert.NoError(t, err)

			proof, err := scram.ClientProof()
			if err != nil {
				if tt.wantErr == nil {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
				}
			} else {
				assert.NotEmpty(t, proof)
				proofStr := hex.EncodeToString(proof)
				assert.Equal(t, tt.wantClientProof, proofStr)
				t.Logf("\n%s", scram.String())
			}
		})
	}
}
