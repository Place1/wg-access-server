package storage

import (
	"errors"
	"strings"
)

var memory = map[string]*Device{}

// implements Storage interface
type InMemoryStorage struct{}

func NewMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Save(device *Device) error {
	memory[key(device)] = device
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
	for key, device := range memory {
		if strings.HasPrefix(key, prefix) {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (s *InMemoryStorage) Get(owner string, name string) (*Device, error) {
	device, ok := memory[keyStr(owner, name)]
	if !ok {
		return nil, errors.New("device doesn't exist")
	}
	return device, nil
}

func (s *InMemoryStorage) Delete(device *Device) error {
	delete(memory, key(device))
	return nil
}
