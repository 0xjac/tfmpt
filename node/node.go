package node

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var _ Node = (Leaf)(nil)
var _ Node = (Hashed)(nil)

var Nil = Leaf([]byte{})

type Node interface {
	ComputeHash() (compact Node, hash Node)
	Hash() Node
}

type Leaf []byte

func (l Leaf) ComputeHash() (Node, Node) { return l, l }

func (l Leaf) Hash() Node { return l }

type Hashed []byte

func (h Hashed) ComputeHash() (Node, Node) { return h, h }

func (h Hashed) Hash() Node { return h }
