package storage

import (
	"time"
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
