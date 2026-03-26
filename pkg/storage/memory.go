package storage

import "sync"

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
	m.Lock()
	defer m.Unlock()

	m.data[key] = value
	return nil
}

func (m *MemoryStore) Update(key string, value []byte) error {
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

	delete(m.data, key)
	return nil
}
