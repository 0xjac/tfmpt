// Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved

package encoding

const (
	AlphabetSize = 1 << nibbleSize // AlphabetSize represent the 16 possible nibble values.
	terminator   = 0x10            // The terminator is the last byte of the path.
	flag         = 0b1             // Arbitrary flag value to set.
	oddPos       = 0               // oddPos is the bit-position of the odd flag.
	oddFlag      = flag << oddPos  // oddFlag bit set at the right position.
	termPos      = 1               // termPos is the bit-position of the terminator flag.
	termFlag     = flag << termPos // termFlag bit set at the right position.
	nibbleSize   = 4               // The nibbleSize is 4 bits or half a byte (convention).
)

// ToHex encodes a key into a byte sequence of hex-encoded nibbles with the 0x10 terminator nibble.
func ToHex(key []byte) []byte {
	nibbles := make([]byte, 0, len(key)*2+1)
	for _, k := range key {
		nibbles = append(nibbles, k>>nibbleSize, k%AlphabetSize)
	}

	nibbles = append(nibbles, terminator)

	return nibbles
}

// CommonPrefixLen returns the length of the common prefix between to paths.
func CommonPrefixLen[T comparable](pathA, pathB []T) int {
	var i int

	for i = 0; i < min(len(pathA), len(pathB)); i++ {
		if pathA[i] != pathB[i] {
			break
		}
	}

	return i
}

// Compact returns a hex-encoded key compacted in a byte sequence where each byte holds two nibbles.
func Compact(hex []byte) []byte {
	var term byte
	if HexKeyHasTerm(hex) {
		term = flag
		hex = hex[:len(hex)-1]
	} else {
		term = 0 // No term flag
	}

	buf := make([]byte, len(hex)/2+1)       // Prefix requires +1.
	buf[0] = term << (nibbleSize + termPos) // Set terminator flag.

	if len(hex)&1 == 1 { // Check if length is odd.
		buf[0] |= flag << (nibbleSize + oddPos) // Set the odd flag.
		buf[0] |= hex[0]                        // Put (XOR) the first nibble in the first byte.
		hex = hex[1:]                           // Remove first nibble from the hex key.
	}

	// Compress nibbles into bytes.
	// Where bi starts at 1 to exclude the prefix taking the first slot.
	// Where ni increments 2 by 2 as each byte consumes two nibbles.
	for bi, ni := 1, 0; ni < len(hex); bi, ni = bi+1, ni+2 {
		buf[bi] = hex[ni]<<nibbleSize | hex[ni+1]
	}

	return buf
}

// ExpandToHex encodes a key to hex-encoded nibbles, removing the prefix and terminator if present.
func ExpandToHex(compact []byte) []byte {
	if len(compact) == 0 {
		return compact
	}
	hex := ToHex(compact)

	// Check if the terminator flag is present.
	// This is a trick: the terminator flag is the second bit and highest flag that can be set on
	// the nibble. Hence, the value with a terminator flag is at least 0b10 (=2).
	// And any value without the terminator flag set is less than 2.
	if hex[0] < termFlag {
		hex = hex[:len(hex)-1] // Remove the terminator if the terminator flag is absent.
	}

	// If the odd flag is set, only remove the prefix in the first nibble.
	// The second nibble is the first actual nibble of the path.
	if (hex[0] & oddFlag) == flag {
		return hex[1:]
	}

	// Otherwise key is even, remove the first two nibbles: the prefix and the empty nibble.
	return hex[2:]
}

// HexKeyHasTerm indicates whether a hex-encoded key ends with a terminator (0x10) byte.
func HexKeyHasTerm(key []byte) bool {
	return len(key) > 0 && key[len(key)-1] == terminator
}
