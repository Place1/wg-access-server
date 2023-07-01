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
	Pingable
	Save(device *Device) error
	List(owner string) ([]*Device, error)
	Get(owner string, name string) (*Device, error)
	GetByPublicKey(publicKey string) (*Device, error)
	Delete(device *Device) error
	Close() error
	Open() error
}

type Watcher interface {
	OnAdd(cb Callback)
	OnDelete(cb Callback)
	OnReconnect(func())
	EmitAdd(device *Device)
	EmitDelete(device *Device)
}

type Pingable interface {
	Ping() error
}

type Callback func(device *Device)

type Device struct {
	Owner         string    `json:"owner" gorm:"type:varchar(100);unique_index:key;primary_key"`
	OwnerName     string    `json:"owner_name"`
	OwnerEmail    string    `json:"owner_email"`
	OwnerProvider string    `json:"owner_provider"`
	Name          string    `json:"name" gorm:"type:varchar(100);unique_index:key;primary_key"`
	PublicKey     string    `json:"public_key" gorm:"unique_index"`
	PresharedKey  string    `json:"preshared_key" gorm:"type:varchar(100)"`
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
		return nil, errors.Wrap(err, "Error parsing storage uri")
	}

	switch u.Scheme {
	case "memory":
		logrus.Warn("Storing data in memory - devices will not persist between restarts")
		return NewMemoryStorage(), nil
	case "postgresql":
		fallthrough
	case "postgres":
		fallthrough
	case "mysql":
		fallthrough
	case "sqlite3":
		logrus.Infof("Storing data in SQL backend at %s", u)
		return NewSqlStorage(u), nil
	}

	return nil, fmt.Errorf("Unknown storage backend %s", u)
}
