package storage

import (
	"errors"
	"strings"
)

var (
	ErrNotFound   = errors.New("key not found")
	ErrInvalidKey = errors.New("key cannot end with a slash")
)

type Store interface {
	Put(key string, value []byte) error
	Update(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	GetBucket(prefix string) (map[string][]byte, error)
}

func IsBucket(key string) bool {
	return strings.HasSuffix(key, "/")
}
