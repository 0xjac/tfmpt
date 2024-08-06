// Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved

// Package crypto shamelessly taken from github.com/ethereum/go-ethereum/crypto to avoid the warning
// related to the import of golang.org/x/crypto v0.22.0 which contains the vulnerability
// CVE-2023-42818.
package crypto

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(KeccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}
