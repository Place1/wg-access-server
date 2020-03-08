package storage

import (
	"time"
)

type Storage interface {
	Save(key string, device *Device) error
	List(prefix string) ([]*Device, error)
	Get(key string) (*Device, error)
	Delete(key string) error
}

type Device struct {
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	PublicKey string    `json:"publicKey"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`

	/**
	 * Metadata fields below.
	 * All metadata tracking can be disabled
	 * from the config file.
	 */

	// metadata about the device during the current session
	LastHandshakeTime time.Time `json:"lastHandshakeTime"`
	ReceiveBytes      int64     `json:"receivedBytes"`
	TransmitBytes     int64     `json:"transmitBytes"`
	Endpoint          string    `json:"endpoint"`

	// metadata about the device through it's lifetime
	LifetimeReceivedBytes int64 `json:"lifetimeReceivedBytes"`
	LifetimeTransmitBytes int64 `json:"lifetimeTransmitBytes"`
}
