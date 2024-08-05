package tfmpt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"

	"tfmpt.local/store"
)

func TestTrieGet(t *testing.T) {
	for _, commit := range []bool{false, true} {
		for _, test := range nodes {
			t.Run(
				fmt.Sprintf("Get[k=%s,v=%s]%s", test.key, test.val, suffix(t, commit)),
				func(t *testing.T) {
					t.Parallel()

					mpt, cleanup := trieSetup(t, commit)

					actual, err := mpt.Get(test.key)

					assertPresent(t, test.key, actual, test.val, err)

					cleanup()
				})
		}

		for _, test := range missings {
			t.Run(fmt.Sprintf("Get[missing_k=%s]", test), func(t *testing.T) {
				t.Parallel()

				mpt, cleanup := trieSetup(t, commit)

				actual, err := mpt.Get(test)

				assertMissing(t, test, actual, err)

				cleanup()
			})
		}
	}

}

func TestTrieDelete(t *testing.T) {
	newVal := []byte("<new_val>")

	for _, commit := range []bool{false, true} {
		for _, test := range nodes {
			t.Run(fmt.Sprintf("Del[k=%s]%s", test.key, suffix(t, commit)), func(t *testing.T) {
				t.Parallel()

				mpt, cleanup := trieSetup(t, commit)

				if err := mpt.Del(test.key); err != nil {
					t.Errorf("Expected key=%s to be deleted, got err=%s", test.key, err)
				}

				for _, node := range nodes {
					val, err := mpt.Get(node.key)

					if bytes.Equal(test.key, node.key) {
						assertMissing(t, node.key, val, err)

						mpt.Put(node.key, newVal)

						val, err = mpt.Get(node.key)
						assertPresent(t, node.key, val, newVal, err)

						continue
					}

					assertPresent(t, node.key, val, node.val, err)
				}

				cleanup()
			})
		}
	}
}

func TestTrieDeleteThenCommit(t *testing.T) {
	newVal := []byte("<new_val>")

	for _, test := range nodes {
		t.Run(fmt.Sprintf("Del[k=%s]", test.key), func(t *testing.T) {
			t.Parallel()
			db, cleanup := storageFixture(t)
			mpt := trieFixture(t, db)

			if err := mpt.Del(test.key); err != nil {
				t.Errorf("Expected key=%s to be deleted, got err=%s", test.key, err)
			}

			root := mpt.Commit()
			mpt = LoadTrie(db, root)

			for _, node := range nodes {
				val, err := mpt.Get(node.key)
				if bytes.Equal(test.key, node.key) {
					assertMissing(t, node.key, val, err)

					mpt.Put(node.key, newVal)

					val, err = mpt.Get(node.key)
					assertPresent(t, node.key, val, newVal, err)

					continue
				}

				assertPresent(t, node.key, val, node.val, err)
			}

			cleanup()
		})
	}
}

func TestTrieProof(t *testing.T) {
	for _, commit := range []bool{false, true} {
		for _, test := range nodes {
			t.Run(fmt.Sprintf("Proof[k=%s]%s", test.key, suffix(t, commit)), func(t *testing.T) {
				t.Parallel()
				ethMPT := ethTrieFixture(t)

				mpt, cleanup := trieSetup(t, commit)

				actual, err := mpt.Proof(test.key)
				if err != nil {
					t.Errorf("Expected a proof for valid key=%s, got err=%s", test.key, err)
				}

				expected := newMockEthProofDB(len(actual))

				if err = ethMPT.Prove(test.key, expected); err != nil {
					t.Errorf("Expected a proof for valid key=%s, got err=%s", test.key, err)
				}

				if len(actual) != len(expected) {
					t.Fatalf("Expected proof length=%d, got length=%d", len(expected), len(actual))
				}

				for _, part := range actual {
					if _, ok := expected[string(part)]; !ok {
						t.Errorf("Bad proof part=%064x", part)
					}
				}

				cleanup()
			})
		}

		for _, test := range missings {
			t.Run(
				fmt.Sprintf("Proof[missing_k=%s]%s", test, suffix(t, commit)),
				func(t *testing.T) {
					t.Parallel()

					mpt, cleanup := trieSetup(t, commit)

					proof, err := mpt.Proof(test)
					if !errors.Is(err, ErrNotFound) {
						t.Errorf("Expected key=%s to be missing, got err=%s", test, err)
					}

					if proof != nil {
						t.Errorf("Expected proof for key=%s to be nil, got %064x", test, proof)
					}

					cleanup()
				})
		}
	}
}

func assertPresent(t *testing.T, key, val, expected []byte, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Expected key=%s to be present, got err=%s", key, err)
	}

	if !bytes.Equal(val, expected) {
		t.Errorf("Expected key=%s to be %s, got val=%s", key, expected, val)
	}
}

func assertMissing(t *testing.T, key, val []byte, err error) {
	t.Helper()

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected key=%s to be missing, got err=%s", key, err)
	}

	if val != nil {
		t.Errorf("Expected key=%s to be missing, got val=%s", key, val)
	}
}

var nodes = []struct {
	key []byte
	val []byte
}{
	{key: []byte("do"), val: []byte("verb")},
	{key: []byte("dog"), val: []byte("puppy")},
	{key: []byte("doge"), val: []byte("coins")},
	{key: []byte("horse"), val: []byte("stallion")},
}

var missings = [][]byte{
	[]byte("<missing>"),
	[]byte("d"),
	[]byte("dogs"),
}

func suffix(t *testing.T, commit bool) string {
	t.Helper()

	if !commit {
		return ""
	}
	return "-commit"
}

func trieSetup(t *testing.T, commit bool) (iTrie, func()) {
	var db store.DB

	cleanup := func() {}
	if commit {
		db, cleanup = storageFixture(t)
	}

	mpt := trieFixture(t, db)

	if commit {
		root := mpt.Commit()
		mpt = LoadTrie(db, root)
	}

	return mpt, cleanup
}

func trieFixture(t *testing.T, db store.DB) iTrie {
	t.Helper()

	mpt := NewEmptyTrie(db)
	for _, node := range nodes {
		mpt.Put(node.key, node.val)
	}

	return mpt
}

// ethTrieFixture returns a trie from the official go-ethereum implementation for comparison.
func ethTrieFixture(t *testing.T) *trie.Trie {
	t.Helper()

	mpt := trie.NewEmpty(nil)
	for _, node := range nodes {
		mpt.MustUpdate(node.key, node.val)
	}

	return mpt
}

func storageFixture(t *testing.T) (store.DB, func()) {
	t.Helper()

	dbPath, err := os.MkdirTemp("", "tlmpt-trietest-db-*")
	if err != nil {
		t.Fatal(err)
	}

	db, err := store.NewLevelDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		var cleanupErr error
		if closeErr := db.Close(); closeErr != nil {
			cleanupErr = fmt.Errorf("failed to close db: %w", closeErr)
		}

		if rmError := os.RemoveAll(dbPath); rmError != nil {
			if cleanupErr != nil {
				cleanupErr = fmt.Errorf("%w, failed to remove db: %w", cleanupErr, rmError)
			} else {
				cleanupErr = fmt.Errorf("failed to remove db: %w", rmError)
			}
		}

		if cleanupErr != nil {
			t.Fatal(cleanupErr)
		}
	}

	return db, cleanup
}

var _ ethdb.KeyValueWriter = (mockEthProofDB)(nil)

type mockEthProofDB map[string][]byte

func (p mockEthProofDB) Put(key []byte, value []byte) error {
	p[string(key)] = value

	return nil
}

func (p mockEthProofDB) Delete(key []byte) error {
	delete(p, string(key))

	return nil
}

func newMockEthProofDB(capacity int) mockEthProofDB {
	return make(map[string][]byte, capacity)
}
