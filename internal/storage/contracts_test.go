package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	require := require.New(t)

	s, err := NewStorage("memory://")
	require.NoError(err)

	require.IsType(&InMemoryStorage{}, s)
}

func TestPostgresqlStorage(t *testing.T) {
	require := require.New(t)

	s, err := NewStorage("postgresql://localhost:5432/dbname?sslmode=disable")
	require.NoError(err)

	require.IsType(&SQLStorage{}, s)
}

func TestMysqlStorage(t *testing.T) {
	require := require.New(t)

	s, err := NewStorage("mysql://localhost:1234/dbname?sslmode=disable")
	require.NoError(err)

	require.IsType(&SQLStorage{}, s)
}

func TestSqliteStorage(t *testing.T) {
	require := require.New(t)

	s, err := NewStorage("sqlite3:///some/path/sqlite.db")
	require.NoError(err)

	require.IsType(&SQLStorage{}, s)
}

func TestSqliteStorageRelativePath(t *testing.T) {
	require := require.New(t)

	s, err := NewStorage("sqlite3://sqlite.db")
	require.NoError(err)

	require.IsType(&SQLStorage{}, s)
}

func TestUnknownStorage(t *testing.T) {
	require := require.New(t)

	s, err := NewStorage("foo://")
	require.Nil(s)
	require.Error(err)
	require.Equal(err.Error(), "Unknown storage backend foo:")
}
