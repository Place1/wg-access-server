package storage

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	Watcher
	Save(device *Device) error
	List(owner string) ([]*Device, error)
	Get(owner string, name string) (*Device, error)
	Delete(device *Device) error
	Close() error
	Open() error
}

type Watcher interface {
	OnAdd(cb Callback)
	OnDelete(cb Callback)
}

type Callback func(device *Device)

type Device struct {
	Owner         string    `json:"owner" gorm:"type:varchar(100);unique_index:key;primary_key"`
	OwnerName     string    `json:"owner_name"`
	OwnerEmail    string    `json:"owner_email"`
	OwnerProvider string    `json:"owner_provider"`
	Name          string    `json:"name" gorm:"type:varchar(100);unique_index:key;primary_key"`
	PublicKey     string    `json:"public_key"`
	Address       string    `json:"address"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`

	/**
	 * Metadata fields below.
	 * All metadata tracking can be disabled
	 * from the config file.
	 */

	// metadata about the device during the current session
	LastHandshakeTime *time.Time `json:"last_handshake_time"`
	ReceiveBytes      int64      `json:"received_bytes"`
	TransmitBytes     int64      `json:"transmit_bytes"`
	Endpoint          string     `json:"endpoint"`
}

func NewStorage(uri string) (Storage, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing storage uri")
	}

	switch u.Scheme {
	case "memory":
		logrus.Warn("storing data in memory - devices will not persist between restarts")
		return NewMemoryStorage(), nil
	case "file":
		logrus.Warn("[DEPRECATION NOTICE] using file:// storage is deprecated and will be removed in an upcoming minor release. Please use sqlite3:// for filesystem based storage.")
		logrus.Infof("storing data in %s", u.Path)
		return NewFileStorage(u.Path), nil
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		logrus.Infof("storing data in SQL backend %s", u.Scheme)
		return NewSqlStorage(u), nil
	}

	return nil, fmt.Errorf("unknown storage backend %s", u.Scheme)
}
