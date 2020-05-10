package storage

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	Save(device *Device) error
	List(owner string) ([]*Device, error)
	Get(owner string, name string) (*Device, error)
	Delete(device *Device) error
	Close() error
	Open() error
}

type Device struct {
	Owner         string    `json:"owner" gorm:"type:varchar(100);unique_index:key"`
	OwnerName     string    `json:"ownerName"`
	OwnerEmail    string    `json:"ownerEmail"`
	OwnerProvider string    `json:"ownerProvider"`
	Name          string    `json:"name" gorm:"type:varchar(100);unique_index:key"`
	PublicKey     string    `json:"publicKey"`
	Address       string    `json:"address"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`

	/**
	 * Metadata fields below.
	 * All metadata tracking can be disabled
	 * from the config file.
	 */

	// metadata about the device during the current session
	LastHandshakeTime *time.Time `json:"lastHandshakeTime"`
	ReceiveBytes      int64      `json:"receivedBytes"`
	TransmitBytes     int64      `json:"transmitBytes"`
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
