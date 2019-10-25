package storage

import (
	"time"
)

type Storage interface {
	Save(*Device) error
	List() ([]*Device, error)
	Get(string) (*Device, error)
	Delete(*Device) error
}

type Device struct {
	Name            string    `json:"name"`
	PublicKey       string    `json:"publicKey"`
	Endpoint        string    `json:"endpoint"`
	Address         string    `json:"address"`
	DNS             string    `json:"dns"`
	CreatedAt       time.Time `json:"createdAt"`
	ServerPublicKey string    `json:"serverPublicKey"`
}
