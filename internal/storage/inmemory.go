package storage

import (
	"errors"
	"strings"
)

// implements Storage interface
type InMemoryStorage struct {
	*InProcessWatcher
	db map[string]*Device
}

func NewMemoryStorage() *InMemoryStorage {
	db := make(map[string]*Device)
	return &InMemoryStorage{
		InProcessWatcher: NewInProcessWatcher(),
		db:               db,
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
	s.EmitAdd(device)
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

func (s *InMemoryStorage) GetByPublicKey(publicKey string) (*Device, error) {
	devices, err := s.List("")
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if device.PublicKey == publicKey {
			return device, nil
		}
	}
	return nil, errors.New("device doesn't exist")
}

func (s *InMemoryStorage) Delete(device *Device) error {
	delete(s.db, key(device))
	s.EmitDelete(device)
	return nil
}

func (s *InMemoryStorage) Ping() error {
	return nil
}
