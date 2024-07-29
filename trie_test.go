package tfmpt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"tfmpt.local/store"
)

func TestTrie_Get(t *testing.T) {
	trie := trieFixture(t, nil)

	val, err := trie.Get([]byte("doge"))

	if err != nil {
		t.Fatalf("Expected key 'doge' to be present, got err: %v", err)
	}

	if !bytes.Equal(val, []byte("coins")) {
		t.Fatalf("Expected value to be 'verb', got %v", val)
	}

	val, err = trie.Get([]byte("dogs"))
	if val != nil {
		t.Fatalf("Key 'dogs' should not point to a value, got %v, err=%v", val, err)
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Key 'dogs' should not point to a value: Expected ErrNotFound, got %v", err)
	}
}

func TestTrie_Proof(t *testing.T) {
	trie := trieFixture(t, nil)

	proof, err := trie.Proof([]byte("horse"))

	if err != nil {
		t.Fatalf("Key 'horse' is present a proof shoulf be generated, got err: %v", err)
	}

	t.Logf("proof %v", proof)
}

func TestTrie_Commit(t *testing.T) {
	db, cleanup := storageFixture(t)
	trie := trieFixture(t, db)
	defer cleanup()

	root := trie.Commit()

	trie2 := LoadTrie(db, root)

	val, err := trie2.Get([]byte("doge"))

	if err != nil {
		t.Fatalf("Expected key 'doge' to be present, got err: %v", err)
	}

	if !bytes.Equal(val, []byte("coins")) {
		t.Fatalf("Expected value to be 'verb', got %v", val)
	}

}

func storageFixture(t *testing.T) (store.LevelDB, func()) {
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

func trieFixture(t *testing.T, db store.DB) iTrie {
	t.Helper()

	trie := NewEmptyTrie(db)

	trie.Put([]byte("do"), []byte("verb"))
	trie.Put([]byte("dog"), []byte("puppy"))
	trie.Put([]byte("doge"), []byte("coins"))
	trie.Put([]byte("horse"), []byte("stallion"))

	return trie
}
