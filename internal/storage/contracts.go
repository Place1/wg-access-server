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
}
