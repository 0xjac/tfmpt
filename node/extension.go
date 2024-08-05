package node

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/rlp"

	"tfmpt.local/encoding"
)

var _ Node = (*Extension)(nil)

type Extension struct {
	Key   []byte
	Next  Node
	Cache Hashed
}

func (e *Extension) String() string {
	return fmt.Sprintf("[%X, %s]", e.Key, e.Next)
}

func (e *Extension) Hash() Node {
	if e.Cache != nil {
		return e.Cache
	}

	hashed := e.Copy()

	hashed.Key = encoding.Compact(e.Key)

	switch e.Next.(type) {
	case *Branch, *Extension:
		hashed.Next = e.Next.Hash()
	}

	hash := hashNode(hashed)
	if cache, ok := hash.(Hashed); ok {
		e.Cache = cache
	} else {
		e.Cache = nil
	}

	return hash
}

func (e *Extension) EncodeRLP(w io.Writer) error {
	eb := rlp.NewEncoderBuffer(w)
	offset := eb.List()
	eb.WriteBytes(e.Key)

	if e.Next == nil {
		if _, err := eb.Write(rlp.EmptyString); err != nil {
			return err
		}
	} else {
		if err := rlp.Encode(eb, e.Next); err != nil {
			return err
		}
	}

	eb.ListEnd(offset)

	return nil
}

func (e *Extension) Copy() *Extension {
	deref := *e
	return &deref
}

func NewExtension(key []byte, next Node, cache Hashed) *Extension {
	return &Extension{Key: key, Next: next, Cache: cache}
}
