package storage

import "errors"

// implements Storage interface
type InMemoryStorage struct{}

var memory = map[string]*Device{}

func NewMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

func (s *InMemoryStorage) Save(device *Device) error {
	memory[device.Name] = device
	return nil
}

func (s *InMemoryStorage) List() ([]*Device, error) {
	devices := []*Device{}
	for _, device := range memory {
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *InMemoryStorage) Get(name string) (*Device, error) {
	device, ok := memory[name]
	if !ok {
		return nil, errors.New("device doesn't exist")
	}
	return device, nil
}

func (s *InMemoryStorage) Delete(device *Device) error {
	delete(memory, device.Name)
	return nil
}
