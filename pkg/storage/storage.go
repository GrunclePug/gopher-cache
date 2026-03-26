package storage

import "errors"

var ErrNotFound = errors.New("key not found")

type Store interface {
	Put(key string, value []byte) error
	Update(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
