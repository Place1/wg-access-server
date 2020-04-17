package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// implements Storage interface
type PostgresqlStorage struct {
	connectionString string
}

func (s *PostgresqlStorage) conn() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", s.connectionString)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to connect to postgresql"))
	}
	// db.LogMode(true)
	return db, err
}

func NewPostgresStorage(connectionString string) *PostgresqlStorage {
	s := &PostgresqlStorage{connectionString}
	db, err := s.conn()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Device{})
	return s
}

func (s *PostgresqlStorage) Save(device *Device) error {
	db, err := s.conn()
	defer db.Close()
	if err != nil {
		return err
	}

	logrus.Debugf("saving device %s", key(device))
	if err := db.Save(&device).Error; err != nil {
		return errors.Wrapf(err, "failed to write device")
	}
	return nil
}

func (s *PostgresqlStorage) List(username string) ([]*Device, error) {
	db, err := s.conn()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	devices := []*Device{}
	if username != "" {
		err = db.Where("owner = ?", username).Find(&devices).Error
	} else {
		err = db.Find(&devices).Error
	}

	logrus.Debugf("Found devices: %+v", devices)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read devices from sql")
	}
	return devices, nil
}

func (s *PostgresqlStorage) Get(owner string, name string) (*Device, error) {
	db, err := s.conn()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	device := &Device{}
	if err := db.Where("owner = ? AND name = ?", owner, name).First(&device).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to read device")
	}
	return device, nil
}

func (s *PostgresqlStorage) Delete(device *Device) error {
	db, err := s.conn()
	defer db.Close()

	if err != nil {
		return err
	}

	if err := db.Delete(&device).Error; err != nil {
		return errors.Wrap(err, "failed to delete device file")
	}
	return nil
}
