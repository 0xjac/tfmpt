package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"tfmpt.local/encoding"
)

func hashNode(n Node) Node {
	rlpEnc, err := rlp.EncodeToBytes(n)
	if err != nil {
		panic(err)
	}

	if len(rlpEnc) < 32 {
		// If node is an extension with a small leaf, it should still be hashed.
		if ext, ok := n.(*Extension); ok {
			if _, ok = ext.Next.(Leaf); ok {
				tmp := Hashed(crypto.Keccak256(rlpEnc))
				return tmp
			}
			return ext
		}

		return n
	}

	if e, ok := n.(*Extension); ok {
		if _, ok := e.Next.(Leaf); ok {
			return Hashed(crypto.Keccak256(rlpEnc))
		}
	}

	return Hashed(crypto.Keccak256(rlpEnc))
}

func Decode(raw []byte, hashed Hashed) (Node, error) {
	items, _, err := rlp.SplitList(raw)
	if err != nil {
		return nil, err
	}

	l, err := rlp.CountValues(items)
	if err != nil {
		return nil, err
	}

	switch l {
	case 2:
		compactKey, rest, err := rlp.SplitString(items)
		if err != nil {
			return nil, err
		}

		key := encoding.ExpandToHex(compactKey)
		if encoding.HexKeyHasTerm(key) { // Leaf node
			data, _, err := rlp.SplitString(rest)
			if err != nil {
				return nil, err
			}

			ext := NewExtension(key, Leaf(data))

			ext.hash = hashed

			return ext, nil
		}

		var next Node
		if next, _, err = decodeHashedChild(rest); err != nil {
			return nil, err
		}

		ext := NewExtension(key, next)

		ext.hash = hashed

		return ext, nil

	case BranchSize:
		b := NewBranch()

		for i := 0; i < BranchChildren; i++ {
			child, rest, err := decodeHashedChild(items)
			if err != nil {
				return nil, err
			}

			b.Children[i], items = child, rest
		}

		value, _, err := rlp.SplitString(items)
		if err != nil {
			return nil, err
		}

		if len(value) > 0 {
			b.Children[BranchValue] = Leaf(value)
		}

		return b, nil

	default:
		return nil, fmt.Errorf("invalid number of items in list: %v", l)
	}
}

func decodeHashedChild(raw []byte) (Node, []byte, error) {
	kind, data, rest, err := rlp.Split(raw)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case kind == rlp.List:
		child, err := Decode(raw, nil)
		return child, rest, err

	case kind == rlp.String && len(data) == 0: // Empty node
		return nil, rest, nil

	case kind == rlp.String && len(data) == 32: // Hash node
		return Hashed(data), rest, nil

	case kind == rlp.String:
		return nil, nil, fmt.Errorf("bad string size %d, expected %d or %d", len(data), 0, 32)

	default:
		return nil, nil, fmt.Errorf("bad rlp kind: %v", kind)
	}
}
