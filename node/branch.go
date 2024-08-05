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
	Cache    Hashed
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
	if b.Cache != nil {
		return b.Cache
	}

	hashed := b.Copy()

	for i := 0; i < BranchChildren; i++ {
		if child := b.Children[i]; child != nil {
			hashed.Children[i] = child.Hash()
		} else {
			hashed.Children[i] = nil
		}
	}

	hash := hashNode(hashed)
	if cache, ok := hash.(Hashed); ok {
		b.Cache = cache
	} else {
		b.Cache = nil
	}

	return hash
}

func (b *Branch) Copy() *Branch {
	deref := *b
	return &deref
}

func NewBranch(cache Hashed) *Branch {
	return &Branch{Children: [BranchSize]Node{}, Cache: cache}
}
