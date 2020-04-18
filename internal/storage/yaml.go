package storage

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StorageWrapper struct {
	Storage
}

func (storageDriver *StorageWrapper) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var connectionStr string
	if err := unmarshal(&connectionStr); err != nil {
		return err
	}

	s, err := NewStorageWrapper(connectionStr)
	if err != nil {
		return err
	}
	*storageDriver = *s
	return nil
}

func NewStorageWrapper(connectionStr string) (*StorageWrapper, error) {
	if connectionStr == "" {
		return &StorageWrapper{NewMemoryStorage()}, nil
	}

	parsedURL, err := url.Parse(connectionStr)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing storage url")
	}

	switch parsedURL.Scheme {
	case "":
	case "memory":
		logrus.Infof("storing data in memory")
		return &StorageWrapper{NewMemoryStorage()}, nil
	case "disk":
		logrus.Infof("storing data in %s", parsedURL.Path)
		return &StorageWrapper{NewDiskStorage(parsedURL.Path)}, nil
	case "postgres":
		password, _ := parsedURL.User.Password()
		decodedQuery, err := url.QueryUnescape(parsedURL.RawQuery)
		if err != nil {
			return nil, errors.Wrap(err, "decoding extra flags")
		}

		storage := NewSqlStorage(
			"postgres",
			fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s %s",
				parsedURL.Hostname(),
				parsedURL.Port(),
				parsedURL.User.Username(),
				password,
				strings.TrimLeft(parsedURL.Path, "/"),
				decodedQuery,
			),
		)
		return &StorageWrapper{storage}, nil
	case "mysql":
		password, _ := parsedURL.User.Password()
		storage := NewSqlStorage(
			"mysql",
			fmt.Sprintf(
				"%s:%s@%s/%s?%s",
				parsedURL.User.Username(),
				password,
				parsedURL.Host,
				strings.TrimLeft(parsedURL.Path, "/"),
				parsedURL.RawQuery,
			),
		)
		return &StorageWrapper{storage}, nil
	case "sqlite3":
		storage := NewSqlStorage(
			"sqlite3",
			parsedURL.Path,
		)
		return &StorageWrapper{storage}, nil
	default:
		return nil, fmt.Errorf("Not a known storage type: %s", parsedURL.Scheme)
	}
	return nil, nil
}
