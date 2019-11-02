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
	Owner           string    `json:"owner"`
	Name            string    `json:"name"`
	PublicKey       string    `json:"publicKey"`
	Endpoint        string    `json:"endpoint"`
	Address         string    `json:"address"`
	DNS             string    `json:"dns"`
	CreatedAt       time.Time `json:"createdAt"`
	ServerPublicKey string    `json:"serverPublicKey"`
}
