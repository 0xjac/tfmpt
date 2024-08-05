package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var _ DB = (*LevelDB)(nil)

type LevelDB struct {
	*leveldb.DB
}

func (l *LevelDB) Get(key []byte) ([]byte, error) {
	return l.DB.Get(key, nil)
}

func (l *LevelDB) Put(key, value []byte) error {
	return l.DB.Put(key, value, nil)
}

func (l *LevelDB) Delete(key []byte) error {
	return l.DB.Delete(key, nil)
}

func NewLevelDB(dbPath string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)

	return &LevelDB{db}, err
}
