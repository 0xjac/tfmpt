package store

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"

	"tfmpt.local/encoding"
	"tfmpt.local/node"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Close() error
}

type Store struct {
	db DB
}

func (s *Store) Commit(n node.Node) (node.Hashed, error) {
	if hashed, err := s.commit(nil, n); err != nil {
		return nil, err
	} else if h, ok := hashed.(node.Hashed); !ok {
		return nil, fmt.Errorf("store: expected Hashed node, got: %T", hashed)
	} else {
		return h, nil
	}
}

func (s *Store) commit(path []byte, n node.Node) (node.Node, error) {
	var err error

	switch current := n.(type) {
	case *node.Branch:
		var nodes [len(current.Children)]node.Node

		for i, child := range current.Children {
			if child == nil {
				continue
			}

			if hashed, ok := child.(node.Hashed); ok {
				nodes[i] = hashed
				continue
			}

			nodes[i], err = s.commit(append(path, byte(i)), child)
		}

		// Branch may hold a "Leaf" which must be explicitly included.
		if current.Children[node.BranchValue] != nil {
			nodes[node.BranchValue] = current.Children[node.BranchValue]
		}

		hashedBranch := current.Copy()
		hashedBranch.Children = nodes

		var encBranch []byte
		if encBranch, err = rlp.EncodeToBytes(hashedBranch); err != nil {
			return nil, err
		}

		err = s.db.Put(path, encBranch)
		if hash := hashedBranch.Hash(); hash != nil {
			return hash, err
		}

		return hashedBranch, err

	case *node.Extension:
		next := current.Next

		if _, ok := next.(*node.Branch); ok {
			next, err = s.commit(append(path, current.Key...), next)
		}

		hashedExtension := current.Copy()
		hashedExtension.Key = encoding.Compact(hashedExtension.Key)

		var encExtension []byte
		if encExtension, err = rlp.EncodeToBytes(hashedExtension); err != nil {
			return nil, err
		}

		err = s.db.Put(path, encExtension)

		if hash := hashedExtension.Hash(); hash != nil {
			return hash, err
		}

		// if _, ok := hashedExtension.Next.(node.Leaf); ok {
		//
		// }

		return hashedExtension, err

	case node.Hashed:
		return current, nil

	case node.Leaf:
		return nil, fmt.Errorf("leaf should not be stored directly")

	default:
		return nil, fmt.Errorf("unknown node type: %T", current)
	}
}

func New(db DB) *Store {
	return &Store{db: db}
}

func (s *Store) Get(path []byte, hash node.Hashed) (node.Node, error) {
	raw, err := s.db.Get(path)
	switch {
	case err != nil:
		return nil, err
	case raw == nil:
		return nil, fmt.Errorf("store: not found")
	}

	n, err := node.Decode(raw, hash)
	if err != nil {
		return nil, fmt.Errorf("store: decode error: %v", err)
	}

	return n, nil

}
