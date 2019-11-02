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

func (s *InMemoryStorage) Save(key string, device *Device) error {
	memory[key] = device
	return nil
}

func (s *InMemoryStorage) List(prefix string) ([]*Device, error) {
	devices := []*Device{}
	for key, device := range memory {
		if strings.HasPrefix(key, prefix) {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (s *InMemoryStorage) Get(key string) (*Device, error) {
	device, ok := memory[key]
	if !ok {
		return nil, errors.New("device doesn't exist")
	}
	return device, nil
}

func (s *InMemoryStorage) Delete(key string) error {
	delete(memory, key)
	return nil
}
