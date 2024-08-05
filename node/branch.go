package node

import (
	"io"

	"github.com/ethereum/go-ethereum/rlp"

	"tfmpt.local/encoding"
)

const (
	BranchChildren = encoding.AlphabetSize
	BranchSize     = encoding.AlphabetSize + 1
	BranchValue    = encoding.AlphabetSize
)

var _ Node = (*Branch)(nil)

type Branch struct {
	Children [BranchSize]Node
	hash     Hashed // hash is a cache for the hash value or nil.
}

func (b *Branch) ComputeHash() (Node, Node) {
	if b.hash != nil {
		return b, b.hash
	}

	compact := b.Copy()
	cachedBranch := b.Copy()
	for i := 0; i < encoding.AlphabetSize; i++ {
		if b.Children[i] == nil {
			compact.Children[i] = Nil
		} else {
			compact.Children[i], cachedBranch.Children[i] = b.Children[i].ComputeHash()
		}
	}

	hashed := hashNode(cachedBranch)
	if h, ok := hashed.(Hashed); ok {
		cachedBranch.hash = h
	}

	return compact, cachedBranch
}

func (b *Branch) Hash() Node {
	// Hash the full node's children, caching the newly hashed subtrees
	hashed := b.Copy()

	for i := 0; i < BranchChildren; i++ {
		if child := b.Children[i]; child != nil {
			hashed.Children[i] = child.Hash()
		} else {
			hashed.Children[i] = nil
		}
	}

	return hash(hashed)
	}

	return h
}

func (b *Branch) Copy() *Branch {
	deref := *b
	return &deref
}

func NewBranch() *Branch { return &Branch{Children: [BranchSize]Node{}, hash: nil} }
