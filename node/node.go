// Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved

package node

var _ Node = (Leaf)(nil)
var _ Node = (Hashed)(nil)

type Node interface {
	Hash() Node
}

type Leaf []byte

func (l Leaf) Hash() Node { return l }

type Hashed []byte

func (h Hashed) Hash() Node { return h }
