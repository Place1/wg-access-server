package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// implements Storage interface
type SQLStorage struct {
	db               *gorm.DB
	sqlType          string
	connectionString string
}

func NewSqlStorage(sqlType string, connectionString string) *SQLStorage {
	return &SQLStorage{
		db:               nil,
		sqlType:          sqlType,
		connectionString: connectionString,
	}
}

func (s *SQLStorage) Open() error {
	db, err := gorm.Open(s.sqlType, s.connectionString)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to connect to %s", s.sqlType))
	}
	s.db = db
	// s.db.LogMode(true)

	// Migrate the schema
	s.db.AutoMigrate(&Device{})
	return nil
}

func (s *SQLStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLStorage) Save(device *Device) error {
	logrus.Debugf("saving device %s", key(device))
	if err := s.db.Save(&device).Error; err != nil {
		return errors.Wrapf(err, "failed to write device")
	}
	return nil
}

func (s *SQLStorage) List(username string) ([]*Device, error) {
	var err error
	devices := []*Device{}
	if username != "" {
		err = s.db.Where("owner = ?", username).Find(&devices).Error
	} else {
		err = s.db.Find(&devices).Error
	}

	logrus.Debugf("Found devices: %+v", devices)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read devices from sql")
	}
	return devices, nil
}

func (s *SQLStorage) Get(owner string, name string) (*Device, error) {
	device := &Device{}
	if err := s.db.Where("owner = ? AND name = ?", owner, name).First(&device).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to read device")
	}
	return device, nil
}

func (s *SQLStorage) Delete(device *Device) error {
	if err := s.db.Delete(&device).Error; err != nil {
		return errors.Wrap(err, "failed to delete device file")
	}
	return nil
}
