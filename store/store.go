// Copyright (C) 2024 Jacques Dafflon | 0xjac - All Rights Reserved

package store

type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Close() error
}
