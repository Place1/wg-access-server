package storage

import (
	"errors"
	"strings"
)

// implements Storage interface
type InMemoryStorage struct {
	db map[string]*Device
}

func NewMemoryStorage() *InMemoryStorage {
	db := make(map[string]*Device)
	return &InMemoryStorage{
		db: db,
	}
}

func (s *InMemoryStorage) Open() error {
	return nil
}

func (s *InMemoryStorage) Close() error {
	return nil
}

func (s *InMemoryStorage) Save(device *Device) error {
	s.db[key(device)] = device
	return nil
}

func (s *InMemoryStorage) List(username string) ([]*Device, error) {
	devices := []*Device{}
	prefix := func() string {
		if username != "" {
			return keyStr(username, "")
		}
		return ""
	}()
	for key, device := range s.db {
		if strings.HasPrefix(key, prefix) {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (s *InMemoryStorage) Get(owner string, name string) (*Device, error) {
	device, ok := s.db[keyStr(owner, name)]
	if !ok {
		return nil, errors.New("device doesn't exist")
	}
	return device, nil
}

func (s *InMemoryStorage) Delete(device *Device) error {
	delete(s.db, key(device))
	return nil
}
