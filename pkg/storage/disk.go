package storage

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type DiskStore struct {
	sync.RWMutex
	dir string
}

func NewDiskStore(dir string) (*DiskStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &DiskStore{
		dir: dir,
	}, nil
}

func (d *DiskStore) path(key string) string {
	return filepath.Join(d.dir, key)
}

func (d *DiskStore) put(key string, value []byte) error {
	return os.WriteFile(d.path(key), value, 0o644)
}

func (d *DiskStore) Put(key string, value []byte) error {
	d.Lock()
	defer d.Unlock()

	return d.put(key, value)
}

func (d *DiskStore) Update(key string, value []byte) error {
	d.Lock()
	defer d.Unlock()

	if _, err := os.Stat(d.path(key)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return err
	}
	return d.put(key, value)
}

func (d *DiskStore) Get(key string) ([]byte, error) {
	d.RLock()
	defer d.RUnlock()

	data, err := os.ReadFile(d.path(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (d *DiskStore) Delete(key string) error {
	d.Lock()
	defer d.Unlock()

	err := os.Remove(d.path(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
