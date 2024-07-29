package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var _ DB = (*LevelDB)(nil)

type LevelDB struct {
	*leveldb.DB
}

func (a LevelDB) Get(key []byte) ([]byte, error) {
	return a.DB.Get(key, nil)
}

func (a LevelDB) Put(key, value []byte) error {
	return a.DB.Put(key, value, nil)
}

func (a LevelDB) Delete(key []byte) error {
	return a.DB.Delete(key, nil)
}

func NewLevelDB(dbPath string) (LevelDB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)

	return LevelDB{db}, err
}
