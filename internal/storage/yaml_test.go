package storage

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type testYaml struct {
	Storage StorageWrapper `yaml:"storage"`
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

func initPostgres() net.Listener {
	// Listen for incoming connections.
	postgresServer, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	go func() {
		for {
			// Listen for an incoming connection.
			conn, err := postgresServer.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				// os.Exit(1)
				return
			}
			go func(conn net.Conn) {
				// Make a buffer to hold incoming data.
				buf := make([]byte, 1024)
				// Read the incoming connection into the buffer.
				_, err := conn.Read(buf)
				if err != nil {
					fmt.Println("Error reading:", err.Error())
				}
				// Send a response back to person contacting us.
				conn.Write([]byte("Message received."))
				// Close the connection when you're done with it.
				conn.Close()
			}(conn)
		}
	}()
	return postgresServer
}

func TestPostgresqlStorage(t *testing.T) {
	var config testYaml

	postgresServer := initPostgres()
	// defer postgresServer.Close()

	yamlStr := fmt.Sprintf("storage: \"postgres://localhost:%d/dbname?sslmode=disable\"", postgresServer.Addr().(*net.TCPAddr).Port)

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

	yamlStr := fmt.Sprintf("storage: \"mysql://localhost:1234/dbname?sslmode=disable\"")

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
