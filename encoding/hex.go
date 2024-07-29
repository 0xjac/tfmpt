package encoding

const AlphabetSize = 16

func ToHex(key []byte) []byte {
	nibbles := make([]byte, 0, len(key)*2+1)
	for _, k := range key {
		nibbles = append(nibbles, k/AlphabetSize, k%AlphabetSize)
	}

	nibbles = append(nibbles, AlphabetSize) // Branch Value is the last child

	return nibbles
}

func PrefixLen(pathA, pathB []byte) int {
	var i int

	for i = 0; i < min(len(pathA), len(pathB)); i++ {
		if pathA[i] != pathB[i] {
			break
		}
	}

	return i
}

func Compact(hex []byte) []byte {
	hexlen := len(hex)
	var term byte
	if KeyHasTerm(hex) {
		term = 1
		hex = hex[:hexlen-1]
		hexlen -= 1
	} else {
		term = 0
	}

	buf := make([]byte, hexlen/2+1) // Prefix requires +1.
	buf[0] = term << 5              // the flag byte
	if hexlen&1 == 1 {
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= hex[0] // first nibble is contained in the first byte
		hex = hex[1:]
	}

	// Nibbles to bytes (bi starts at 1 because of prefix).
	for bi, ni := 1, 0; ni < len(hex); bi, ni = bi+1, ni+2 {
		buf[bi] = hex[ni]<<4 | hex[ni+1]
	}

	return buf
}

func ExpandToHex(compact []byte) []byte {
	if len(compact) == 0 {
		return compact
	}
	hex := ToHex(compact)

	if hex[0] < 2 { // delete terminator flag
		hex = hex[:len(hex)-1]
	}

	cut := 2 - hex[0]&1 // ut from odd flag
	return hex[cut:]
}

func KeyHasTerm(key []byte) bool {
	return len(key) > 0 && key[len(key)-1] == AlphabetSize
}
