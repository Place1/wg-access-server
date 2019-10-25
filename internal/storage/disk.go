package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// implements Storage interface
type DiskStorage struct {
	directory string
}

func NewDiskStorage(directory string) *DiskStorage {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0600); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to create storage directory"))
		}
	}
	return &DiskStorage{directory}
}

func (s *DiskStorage) Save(device *Device) error {
	path := s.deviceFilePath(device.Name)
	logrus.Infof("saving new device %s", path)
	bytes, err := json.Marshal(device)
	if err != nil {
		return errors.Wrap(err, "failed to marshal device")
	}
	if err := ioutil.WriteFile(path, bytes, 0600); err != nil {
		return errors.Wrapf(err, "failed to write device to file %s", path)
	}
	return nil
}

func (s *DiskStorage) List() ([]*Device, error) {
	devices := []*Device{}
	files, err := ioutil.ReadDir(s.directory)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list storage directory")
	}
	for _, file := range files {
		device, err := s.Get(filepath.Base(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))))
		if err != nil {
			return nil, errors.Wrap(err, "failed to read device file")
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *DiskStorage) Get(name string) (*Device, error) {
	path := s.deviceFilePath(name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read device file %s", path)
	}
	device := &Device{}
	if err := json.Unmarshal(bytes, device); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal device file %s", path)
	}
	return device, nil
}

func (s *DiskStorage) Delete(device *Device) error {
	if err := os.Remove(s.deviceFilePath(device.Name)); err != nil {
		return errors.Wrap(err, "failed to delete device file")
	}
	return nil
}

func (s *DiskStorage) deviceFilePath(name string) string {
	// TODO: protect against path traversal
	// and make sure names are reasonably sane
	return filepath.Join(s.directory, fmt.Sprintf("%s.json", name))
}
