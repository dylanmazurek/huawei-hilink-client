package crypto

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name  string
		nonce string
		salt  string
		data  string

		expected string
		wantErr  error
	}{
		{
			name:  "simple decrypt",
			nonce: "3a13f730cb65a370f5986c6ea6cefe3fafa1873f04a00dc26fbddf11254fb3cd",
			salt:  "3a13f730cb65a370f5986c6ea6cefe3fafa1873f04a00dc26fbddf11254fb3cd",
			data:  "c81a754f580eb27bcede74208134f47a8365dfc5fd74cedb90d60bd2637af688a4b5eb4c735f7a9c686d44c9afd30a8622814f27cb6049cd0df3b81be110c3ca085916192996a4c7ce350525150f8ac8e7092bd91049a5e81ee93c8e201510f3f205f38ab407aeb6cc05e61fbb711709e8f989fd1edf0a6e896bfa51c816f0dcb26d7f716098dc1e79aad8153c5591b950f80773a0098e366d7839722986724806f37c17203b482393b32527a41f54df6718df6647800d188116fb502e04d119917248edd1dbaffb9503f36aa2d38771ba91303f6c26b76ed958bc24b6228ec943851f1fe006c79458c90d7c7b74bfb9fd8b8141a67d21b743d1bbd4484c2df0f8951f28f4731fd3be871150eec383e9147f584578cf5233c7963a417508cbb049864cda605652dec9c8e5b018965ed7850f8ff76f7d962618306a9cbc8cc0fd6e42d4cdaeadbd98942295ee056b54a0f7202cd0bd27d2f99d0c39d3ec87e2b99b8d6dfce6540d2a8bdbd7c222f8a8d382489b2c197fd2066c39dcb98501257cf07a313aacc309991b8f97ca7d1795bbf6ed5681fd8b725b15a44334c4662e63f4e032e712b7f32aa5025c2dab2912cd4a43b1d4638648373c47e36bf8dbef15102718967abf5547a1d16ed107e5d0bb4876cdccf0b0a5183a7e74e4ae2b4734086525a52d42618ae696c2713290c4958229ce9d00c141ab2e6a21b0953dd34ce130b52d42a55ee056349e0045a3f2e684476ce01f3fe7288bd7573a0e518d98e6dbeaccda30258832681c39b1776b1798b12314ca6667b11d01ddae27bcb0381938e61dcf021c0c8dd721b29d959a1429ec46f3f187cd1e6fb34c1b2cde386727cf284daec09b0b97d5fb7473bc07ea89da4092d885c59d975245fbc31dd581",

			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scramOpts := []Option{
				WithStaticNonce(tt.nonce, tt.nonce),
			}

			scram, err := NewScram(scramOpts...)
			assert.NoError(t, err)

			decrypted, err := AESDecryptCBC(*scram, tt.nonce, tt.salt, tt.data)
			assert.Equal(t, tt.wantErr, err)

			if decrypted != nil {
				log.Info().Msgf("decrypted: %s", *decrypted)
			}

			assert.Equal(t, tt.expected, decrypted)
		})
	}
}
