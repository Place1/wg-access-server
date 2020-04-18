package storage

import (
	"fmt"
	"log"
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

	s, err := NewStorage(connectionStr)
	if err != nil {
		return err
	}
	*storageDriver = *s
	return nil
}

func NewStorage(connectionStr string) (*StorageWrapper, error) {
	parsedURL, err := url.Parse(connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	switch parsedURL.Scheme {
	case "memory":
		return &StorageWrapper{NewMemoryStorage()}, nil
		break
	case "disk":
		logrus.Infof("storing data in %s", parsedURL.Path)
		return &StorageWrapper{NewDiskStorage(parsedURL.Path)}, nil
		break
	case "postgres":
		password, _ := parsedURL.User.Password()
		decodedQuery, err := url.QueryUnescape(parsedURL.RawQuery)
		if err != nil {
			return nil, errors.Wrap(err, "decoding extra flags")
		}

		storage, err := NewSqlStorage(
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
		if err != nil {
			return nil, errors.Wrap(err, "Connecting to db")
		}
		return &StorageWrapper{storage}, nil
	case "mysql":
		password, _ := parsedURL.User.Password()
		storage, err := NewSqlStorage(
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
		if err != nil {
			return nil, errors.Wrap(err, "Connecting to db")
		}
		return &StorageWrapper{storage}, nil
	case "sqlite3":
		storage, err := NewSqlStorage(
			"sqlite3",
			parsedURL.Path,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Connecting to db")
		}
		return &StorageWrapper{storage}, nil
	default:
		return nil, fmt.Errorf("Not a known storage type: %s", parsedURL.Scheme)
	}
	return nil, nil
}
