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
type FileStorage struct {
	*InProcessWatcher
	directory string
}

func NewFileStorage(directory string) *FileStorage {
	return &FileStorage{
		InProcessWatcher: NewInProcessWatcher(),
		directory:        directory,
	}
}

func (s *FileStorage) Open() error {
	if _, err := os.Stat(s.directory); os.IsNotExist(err) {
		if err := os.MkdirAll(s.directory, 0600); err != nil {
			return errors.Wrap(err, "failed to create storage directory")
		}
	}
	return nil
}

func (s *FileStorage) Close() error {
	return nil
}

func (s *FileStorage) Save(device *Device) error {
	key := key(device)
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
	s.emitAdd(device)
	return nil
}

func (s *FileStorage) List(username string) ([]*Device, error) {
	prefix := func() string {
		if username != "" {
			return keyStr(username, "")
		}
		return ""
	}()
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
		if strings.HasPrefix(p, prefix) && filepath.Ext(path) == ".json" {
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

func (s *FileStorage) Get(owner string, name string) (*Device, error) {
	key := keyStr(owner, name)
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

func (s *FileStorage) Delete(device *Device) error {
	if err := os.Remove(s.deviceFilePath(key(device))); err != nil {
		return errors.Wrap(err, "failed to delete device file")
	}
	s.emitDelete(device)
	return nil
}

func (s *FileStorage) deviceFilePath(key string) string {
	// TODO: protect against path traversal
	// and make sure names are reasonably sane
	return filepath.Join(s.directory, fmt.Sprintf("%s.json", key))
}
