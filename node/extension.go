package node

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/rlp"

	"tfmpt.local/encoding"
)

var _ Node = (*Extension)(nil)

type Extension struct {
	Key  []byte
	Next Node
	hash Hashed // hash is a cache for the hash value or nil.
}

func (e *Extension) String() string {
	return fmt.Sprintf("[%X, %s]", e.Key, e.Next)
}

func (e *Extension) Hash() Node {
	hashed := e.Copy()

	hashed.Key = encoding.Compact(e.Key)

	switch e.Next.(type) {
	case *Branch, *Extension:
		hashed.Next = e.Next.Hash()
	}

	return hash(hashed)
}

func (e *Extension) ComputeHash() (Node, Node) {
	if e.hash != nil {
		return e, e.hash
	}

	compact := e.Copy()
	compact.Key = encoding.Compact(compact.Key)

	cachedExtension := e.Copy()

	switch nxt := e.Next.(type) {
	case *Branch, *Extension:
		compact.Next, cachedExtension.Next = nxt.ComputeHash()
	}

	hashed := hashNode(compact)
	h, ok := hashed.(Hashed)
	if ok {
		cachedExtension.hash = h
	} else {
		cachedExtension.hash = nil
	}

	return compact, cachedExtension
}

func (e *Extension) EncodeRLP(w io.Writer) error {
	eb := rlp.NewEncoderBuffer(w)
	offset := eb.List()
	eb.WriteBytes(e.Key)

	if e.Next == nil {
		eb.Write(rlp.EmptyString)
	} else {
		if encNext, err := rlp.EncodeToBytes(e.Next); err != nil {
			return err
		} else {
			eb.Write(encNext)
		}
	}

	eb.ListEnd(offset)

	return nil
}

func (e *Extension) Copy() *Extension {
	deref := *e
	return &deref
}

func NewExtension(key []byte, next Node) *Extension {
	return &Extension{Key: key, Next: next}
}
