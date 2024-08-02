package tfmpt

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"tfmpt.local/encoding"
	"tfmpt.local/node"
	"tfmpt.local/store"
)

var (
	ErrNotFound = errors.New("not found")

	// emptyRoot is the precomputed hash of an empty MPT.
	// It is equivalent to keccak256(rlp(byte(0)).
	emptyRoot = []byte{
		0x56, 0xE8, 0x1F, 0x17, 0x1B, 0xCC, 0x55, 0xA6,
		0xFF, 0x83, 0x45, 0xE6, 0x92, 0xC0, 0xF8, 0x6E,
		0x5B, 0x48, 0xE0, 0x1B, 0x99, 0x6C, 0xAD, 0xC0,
		0x01, 0x62, 0x2F, 0xB5, 0xE3, 0x63, 0xB4, 0x21}
	_ iTrie = (*Trie)(nil)
)

type iTrie interface {
	// Get returns the value associate with the key
	// error is returned if the key is not found.
	Get(key []byte) ([]byte, error)

	// Put inserts the [key, value] node in the trie
	Put(key []byte, value []byte)

	// Del removes a node from the trie
	// returns an error if not found.
	Del(key []byte) error

	// Commit saves the trie in persistent storage
	// and returns the trie root key.
	Commit() []byte

	// Proof returns the Merkle-proof associated with
	// a node. An error is returned if the node is not found.
	Proof(key []byte) ([][]byte, error)
}

type Trie struct {
	root  node.Node
	store *store.Store
}

func (t *Trie) Get(key []byte) ([]byte, error) {
	path := encoding.ToHex(key)
	return t.get(t.root, path, 0)
}

func (t *Trie) Put(key []byte, value []byte) {
	path := encoding.ToHex(key)
	t.root = t.put(t.root, path, node.Leaf(value))
}

func (t *Trie) Del(key []byte) error {
	// TODO implement me
	panic("implement me")
}

func (t *Trie) Commit() []byte {
	if t.root == nil {
		return emptyRoot
	}

	if t.store == nil {
		panic("store is not set")
	}

	_, t.root = t.root.ComputeHash()

	hashedRoot, err := t.store.Commit(t.root)
	if err != nil {
		panic(err)
	}

	t.root = hashedRoot

	return hashedRoot
}

func (t *Trie) Proof(key []byte) ([][]byte, error) {
	path := encoding.ToHex(key)
	nodes := make([]node.Node, 0, len(path)) // path len is an upper bound on the number of nodes.
	nextNode := t.root

	for len(path) > 0 && nextNode != nil {
		switch current := nextNode.(type) {
		case nil:
			break

		case *node.Branch:
			nextNode = current.Children[key[0]]
			key = key[1:]
			nodes = append(nodes, current)

		case *node.Extension:
			if len(path) < len(current.Key) || !bytes.Equal(current.Key, key[:len(current.Key)]) {
				// The trie doesn't contain the key.
				return nil, ErrNotFound
			} else {
				nextNode = current.Next
				key = key[len(current.Key):]
			}
			nodes = append(nodes, current)

		default:
			return nil, fmt.Errorf("unknown node type: %T", current)
		}
	}

	if len(nodes) == 0 {
		return nil, ErrNotFound
	}

	proof := make([][]byte, 0, len(nodes)) // Nodes len is a safe upper bound.
	for i, n := range nodes {
		var hashedNode node.Node

		n, hashedNode = n.ComputeHash()

		// Only consider the root or nodes with hashes for the proof.
		if hash, ok := hashedNode.(node.Hashed); ok || i == 0 {
			enc, err := rlp.EncodeToBytes(n)
			if err != nil {
				return nil, err
			}

			if !ok {
				hash = crypto.Keccak256(enc)
			}
			proof = append(proof, hash)
		}
	}

	return proof, nil
}

func (t *Trie) get(curr node.Node, path []byte, depth int) ([]byte, error) {
	switch current := curr.(type) {
	case nil:
		return nil, ErrNotFound

	case *node.Branch:
		return t.get(current.Children[path[depth]], path, depth+1)

	case node.Leaf: // Reached end of trie.
		return current, nil
	case *node.Extension:
		keylen := len(current.Key)

		// Match in the middle of extension or the path doesn't match.
		if len(path)-depth < keylen || !bytes.Equal(path[depth:depth+keylen], current.Key) {
			return nil, ErrNotFound
		}

		return t.get(current.Next, path, depth+keylen) // Move through the extension.

	case node.Hashed:
		actual, err := t.store.Get(path[:depth], current)
		if err != nil {
			return nil, err
		}

		return t.get(actual, path, depth)

	default:
		return nil, fmt.Errorf("unknown node type: %T", current)
	}
}

func (t *Trie) put(curr node.Node, path []byte, value node.Node) node.Node {
	if len(path) == 0 { // Trivial we just return the node
		return value
	}

	switch current := curr.(type) {
	case nil:
		return node.NewExtension(path, value)

	case *node.Branch:
		branchKey := path[0]

		current = current.Copy()
		current.Children[branchKey] = t.put(current.Children[branchKey], path[1:], value)

		return current

	case *node.Leaf:
		panic("Leaf should be put with parent Extension")

	case *node.Extension:
		match := encoding.CommonPrefixLen(path, current.Key)
		if match == len(current.Key) { // Path longer than ext, travel down to next node.
			next := t.put(current.Next, path[match:], value)

			return node.NewExtension(current.Key, next)
		}

		// Insert branch after matched prefix.
		branch := node.NewBranch()

		// Insert extension's next as new child.
		branch.Children[current.Key[match]] = t.put(nil, current.Key[match+1:], current.Next)

		// Insert value as new child.
		branch.Children[path[match]] = t.put(nil, path[match+1:], value)

		if match == 0 { // No path before the branch, so no need for an extension.
			return branch
		}

		// Create extension pointing to the branch:
		return node.NewExtension(path[:match], branch)

	default:
		panic(fmt.Sprintf("invalid node type: %T", current))
	}
}

func NewEmptyTrie(db store.DB) *Trie {
	return &Trie{root: nil, store: store.New(db)}
}

func LoadTrie(db store.DB, root node.Hashed) *Trie {
	return &Trie{root: root, store: store.New(db)}
}
