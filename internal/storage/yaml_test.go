package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type testYaml struct {
	Storage StorageWrapper `yaml:"storage"`
}

func TestEmptyStorage(t *testing.T) {
	var config testYaml

	if err := yaml.Unmarshal([]byte("storage: \"\""), &config); err != nil {
		t.Fatal(err, "failed to bind configuration file")
	}

	expected := "*storage.InMemoryStorage"
	actual := fmt.Sprintf("%T", config.Storage.Storage)
	assert.EqualValues(t, string(expected), string(actual))
}

func TestMemoryStorage(t *testing.T) {
	var config testYaml

	if err := yaml.Unmarshal([]byte("storage: \"memory://\""), &config); err != nil {
		t.Fatal(err, "failed to bind configuration file")
	}

	expected := "*storage.InMemoryStorage"
	actual := fmt.Sprintf("%T", config.Storage.Storage)
	assert.EqualValues(t, string(expected), string(actual))
}

func TestDiskStorage(t *testing.T) {
	var config testYaml
	dir, _ := os.Getwd()

	yamlStr := fmt.Sprintf("storage: \"disk://%s\"", filepath.Join(dir, "test"))

	if err := yaml.Unmarshal([]byte(yamlStr), &config); err != nil {
		t.Fatal(err, "failed to bind configuration file")
	}

	expected := "*storage.DiskStorage"
	actual := fmt.Sprintf("%T", config.Storage.Storage)
	assert.EqualValues(t, string(expected), string(actual))
}

func TestPostgresqlStorage(t *testing.T) {
	var config testYaml

	yamlStr := "storage: \"postgres://localhost:5432/dbname?sslmode=disable\""

	if err := yaml.Unmarshal([]byte(yamlStr), &config); err != nil {
		t.Fatal(err, "failed to bind configuration file")
	}
	defer config.Storage.Close()

	expected := "*storage.SQLStorage"
	actual := fmt.Sprintf("%T", config.Storage.Storage)
	assert.EqualValues(t, string(expected), string(actual))
}

func TestMysqlStorage(t *testing.T) {
	var config testYaml

	yamlStr := "storage: \"mysql://localhost:1234/dbname?sslmode=disable\""

	if err := yaml.Unmarshal([]byte(yamlStr), &config); err != nil {
		t.Fatal(err, "failed to bind configuration file")
	}
	defer config.Storage.Close()

	expected := "*storage.SQLStorage"
	actual := fmt.Sprintf("%T", config.Storage.Storage)
	assert.EqualValues(t, string(expected), string(actual))
}

func TestSqliteStorage(t *testing.T) {
	var config testYaml
	dir, _ := os.Getwd()

	yamlStr := fmt.Sprintf("storage: \"sqlite3://%s\"", filepath.Join(dir, "sqlite.db"))

	if err := yaml.Unmarshal([]byte(yamlStr), &config); err != nil {
		t.Fatal(err, "failed to bind configuration file")
	}
	defer config.Storage.Close()

	expected := "*storage.SQLStorage"
	actual := fmt.Sprintf("%T", config.Storage.Storage)
	assert.EqualValues(t, string(expected), string(actual))
}

func TestUnknownStorage(t *testing.T) {
	var config testYaml

	err := yaml.Unmarshal([]byte("storage: \"foo://\""), &config)
	expected := "Not a known storage type: foo"
	actual := err.Error()
	assert.EqualValues(t, string(expected), string(actual))
}
