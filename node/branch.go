// Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved

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

func (b *Branch) EncodeRLP(w io.Writer) error {
	eb := rlp.NewEncoderBuffer(w)
	offset := eb.List()

	for _, child := range &b.Children {
		if child == nil {
			if _, err := eb.Write(rlp.EmptyString); err != nil {
				return err
			}
		} else {
			if err := rlp.Encode(eb, child); err != nil {
				return err
			}
		}
	}

	eb.ListEnd(offset)

	return nil
}

func (b *Branch) Copy() *Branch {
	deref := *b
	return &deref
}

func NewBranch(cache Hashed) *Branch {
	return &Branch{Children: [BranchSize]Node{}, Cache: cache}
}
