package dbfixtures_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDbfixturesSuite(t *testing.T) {
	suite.Run(t, new(dbfixturesUnitTestSuite))
}

/*
mockDriver is a mock of the Driver struct to be used in testing.
*/
type mockDriver struct {
	mock.Mock
}

func (mock *mockDriver) Truncate(tableNames []string) error {
	args := mock.Called(tableNames)

	return args.Error(0)
}
func (mock *mockDriver) InsertFixtures(tableName string, fixtures []interface{}) error {
	args := mock.Called(tableName, fixtures)

	return args.Error(0)
}
func (mock *mockDriver) Close() error {
	args := mock.Called()

	return args.Error(0)
}
