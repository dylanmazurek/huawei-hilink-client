package crypto

import (
	"crypto/sha256"
	"hash"
)

type Options struct {
	hasher       func() hash.Hash
	keySizeBytes int

	password   string
	iterations int

	nonce      *string
	finalNonce *string
}

func DefaultOptions() Options {
	defaultOptions := Options{
		hasher:       sha256.New,
		keySizeBytes: 32,

		iterations: 100,
	}

	return defaultOptions
}

type Option func(*Options)

func WithPassword(password string) Option {
	return func(opts *Options) {
		opts.password = password
	}
}

func WithStaticNonce(nonce string, finalNonce string) Option {
	return func(opts *Options) {
		opts.nonce = &nonce
		opts.finalNonce = &finalNonce
	}
}

func WithHasher(hasher func() hash.Hash) Option {
	return func(opts *Options) {
		opts.hasher = hasher
	}
}

func WithKeySizeBytes(keySizeBytes int) Option {
	return func(opts *Options) {
		opts.keySizeBytes = keySizeBytes
	}
}

func WithIterations(iterations int) Option {
	return func(opts *Options) {
		opts.iterations = iterations
	}
}
