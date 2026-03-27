package storage

import (
	"strings"
	"sync"
)

type MemoryStore struct {
	sync.RWMutex
	data map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

func (m *MemoryStore) Put(key string, value []byte) error {
	if IsBucket(key) {
		return ErrInvalidKey
	}

	m.Lock()
	defer m.Unlock()

	m.data[key] = value
	return nil
}

func (m *MemoryStore) Update(key string, value []byte) error {
	if IsBucket(key) {
		return ErrInvalidKey
	}

	m.Lock()
	defer m.Unlock()

	if _, ok := m.data[key]; !ok {
		return ErrNotFound
	}

	m.data[key] = value
	return nil
}

func (m *MemoryStore) Get(key string) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()

	val, ok := m.data[key]
	if !ok {
		return nil, ErrNotFound
	}
	return val, nil
}

func (m *MemoryStore) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	// Handle Bucket Delete (Recursive)
	if IsBucket(key) {
		for k := range m.data {
			if strings.HasPrefix(k, key) {
				delete(m.data, k)
			}
		}
		return nil
	}

	// Handle Single Key Delete
	if _, ok := m.data[key]; !ok {
		return ErrNotFound
	}
	delete(m.data, key)
	return nil
}

func (m *MemoryStore) GetBucket(prefix string) (map[string][]byte, error) {
	m.RLock()
	defer m.RUnlock()

	results := make(map[string][]byte)
	for k, v := range m.data {
		if strings.HasPrefix(k, prefix) {
			results[k] = v
		}
	}
	return results, nil
}
