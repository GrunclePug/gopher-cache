package storage

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	escapedKey := url.PathEscape(key)
	return filepath.Join(d.dir, escapedKey)
}

func (d *DiskStore) put(key string, value []byte) error {
	return os.WriteFile(d.path(key), value, 0o644)
}

func (d *DiskStore) Put(key string, value []byte) error {
	if IsBucket(key) {
		return ErrInvalidKey
	}

	d.Lock()
	defer d.Unlock()

	return d.put(key, value)
}

func (d *DiskStore) Update(key string, value []byte) error {
	if IsBucket(key) {
		return ErrInvalidKey
	}

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

	if strings.HasSuffix(key, "/") {
		entries, _ := os.ReadDir(d.dir)
		for _, entry := range entries {
			unvKey, _ := url.PathUnescape(entry.Name())
			if strings.HasPrefix(unvKey, key) {
				os.Remove(filepath.Join(d.dir, entry.Name()))
			}
		}
		return nil
	}

	err := os.Remove(d.path(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (d *DiskStore) GetBucket(prefix string) (map[string][]byte, error) {
	d.RLock()
	defer d.RUnlock()

	entries, err := os.ReadDir(d.dir)
	if err != nil {
		return nil, err
	}

	results := make(map[string][]byte)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		key, err := url.PathUnescape(entry.Name())
		if err != nil {
			continue
		}

		if strings.HasPrefix(key, prefix) {
			data, err := os.ReadFile(filepath.Join(d.dir, entry.Name()))
			if err == nil {
				results[key] = data
			}
		}
	}
	return results, nil
}
