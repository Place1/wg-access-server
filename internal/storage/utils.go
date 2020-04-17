package storage

import (
	"path/filepath"
)

func keyStr(owner string, name string) string {
	return filepath.Join(owner, name)
}

func key(device *Device) string {
	return keyStr(device.Owner, device.Name)
}
