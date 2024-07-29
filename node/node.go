package node

var _ Node = (Leaf)(nil)
var _ Node = (Hashed)(nil)

var Nil = Leaf([]byte{})

type Node interface {
	ComputeHash() (compact Node, hash Node)
	Hash() Hashed
}

type Leaf []byte

func (l Leaf) ComputeHash() (Node, Node) { return l, l }

func (l Leaf) Hash() Hashed { return nil }

type Hashed []byte

func (h Hashed) ComputeHash() (Node, Node) { return h, h }

func (h Hashed) Hash() Hashed { return nil }
