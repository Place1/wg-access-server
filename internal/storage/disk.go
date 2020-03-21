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

func (s *DiskStorage) Save(key string, device *Device) error {
	path := s.deviceFilePath(key)
	logrus.Debugf("saving device %s", path)
	bytes, err := json.Marshal(device)
	if err != nil {
		return errors.Wrap(err, "failed to marshal device")
	}
	os.MkdirAll(filepath.Dir(path), 0600)
	if err := ioutil.WriteFile(path, bytes, 0600); err != nil {
		return errors.Wrapf(err, "failed to write device to file %s", path)
	}
	return nil
}

func (s *DiskStorage) List(prefix string) ([]*Device, error) {
	files := []string{}
	err := filepath.Walk(s.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		p := strings.TrimPrefix(path, s.directory)
		p = strings.TrimPrefix(p, string(os.PathSeparator))
		if strings.HasPrefix(p, prefix) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list storage directory")
	}
	logrus.Debugf("Found files: %+v", files)

	devices := []*Device{}
	for _, file := range files {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read device file %s", file)
		}
		device := &Device{}
		if err := json.Unmarshal(bytes, device); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal device file %s", file)
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to read device file")
		}
		devices = append(devices, device)
	}
	logrus.Debugf("Found devices: %+v", devices)
	return devices, nil
}

func (s *DiskStorage) Get(key string) (*Device, error) {
	path := s.deviceFilePath(key)
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

func (s *DiskStorage) Delete(key string) error {
	if err := os.Remove(s.deviceFilePath(key)); err != nil {
		return errors.Wrap(err, "failed to delete device file")
	}
	return nil
}

func (s *DiskStorage) deviceFilePath(key string) string {
	// TODO: protect against path traversal
	// and make sure names are reasonably sane
	return filepath.Join(s.directory, fmt.Sprintf("%s.json", key))
}
