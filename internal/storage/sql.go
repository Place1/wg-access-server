package storage

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// GormLogger is a custom logger for Gorm, making it use logrus.
type GormLogger struct{}

// Print handles log events from Gorm for the custom logger.
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		logrus.WithFields(
			logrus.Fields{
				"module":  "gorm",
				"type":    "sql",
				"rows":    v[5],
				"src_ref": v[1],
				"values":  v[4],
			},
		).Debug(v[3])
	case "logrus":
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "logrus"}).Print(v[2])
	}
}

// implements Storage interface
type SQLStorage struct {
	Watcher
	db               *gorm.DB
	sqlType          string
	connectionString string
}

func NewSqlStorage(u *url.URL) *SQLStorage {
	var connectionString string

	switch u.Scheme {
	case "postgres":
		connectionString = pgconn(u)
	case "mysql":
		connectionString = mysqlconn(u)
	case "sqlite3":
		connectionString = sqlite3conn(u)
	default:
		// unreachable because our storage backend factory
		// function (contracts.go) already checks the url scheme.
		logrus.Panicf("unknown sql storage backend %s", u.Scheme)
	}

	return &SQLStorage{
		Watcher:          nil,
		db:               nil,
		sqlType:          u.Scheme,
		connectionString: connectionString,
	}
}

func pgconn(u *url.URL) string {
	password, _ := u.User.Password()
	decodedQuery, err := url.QueryUnescape(u.RawQuery)
	if err != nil {
		logrus.Warnf("failed to unescape connection string query parameters - they will be ignored")
		decodedQuery = ""
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		u.Hostname(),
		u.Port(),
		u.User.Username(),
		password,
		strings.TrimLeft(u.Path, "/"),
		decodedQuery,
	)
}

func mysqlconn(u *url.URL) string {
	password, _ := u.User.Password()
	return fmt.Sprintf(
		"%s:%s@%s/%s?%s",
		u.User.Username(),
		password,
		u.Host,
		strings.TrimLeft(u.Path, "/"),
		u.RawQuery,
	)
}

func sqlite3conn(u *url.URL) string {
	return filepath.Join(u.Host, u.Path)
}

func (s *SQLStorage) Open() error {
	db, err := gorm.Open(s.sqlType, s.connectionString)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to connect to %s", s.sqlType))
	}
	s.db = db

	db.SetLogger(&GormLogger{})
	db.LogMode(true)

	// Migrate the schema
	s.db.AutoMigrate(&Device{})

	if s.sqlType == "postgres" {
		watcher, err := NewPgWatcher(s.connectionString, db.NewScope(&Device{}).TableName())
		if err != nil {
			return errors.Wrap(err, "failed to create pg watcher")
		}
		s.Watcher = watcher
	} else {
		s.Watcher = NewInProcessWatcher()
	}

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
