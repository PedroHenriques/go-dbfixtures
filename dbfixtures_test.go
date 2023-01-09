package dbfixtures_test

import (
	"testing"

	"github.com/PedroHenriques/go-dbfixtures"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type dbfixturesTestSuite struct {
	suite.Suite
	mockDrivers []*mockDriver
}

func (suite *dbfixturesTestSuite) SetupTest() {
	driver1 := &mockDriver{}
	driver2 := &mockDriver{}

	driver1.truncateDefaultCall = driver1.On("Truncate", mock.Anything).Return(nil)
	driver2.truncateDefaultCall = driver2.On("Truncate", mock.Anything).Return(nil)
	driver1.insertFixturesDefaultCall = driver1.On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)
	driver2.insertFixturesDefaultCall = driver2.On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)

	suite.mockDrivers = []*mockDriver{driver1, driver2}
}

func (suite *dbfixturesTestSuite) TestNewItShouldReturnAnInstanceOfTheTypeDbfixtures() {
	res := dbfixtures.New()

	require.Implements(suite.T(), (*dbfixtures.IDbfixtures)(nil), res)
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverTruncateMethodOnce() {
	fixtures := make(map[string][]interface{})
	fixtures["table1"] = []interface{}{"some object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"table1"}, fixtures)

	suite.mockDrivers[0].AssertNumberOfCalls(suite.T(), "Truncate", 1)
	suite.mockDrivers[1].AssertNumberOfCalls(suite.T(), "Truncate", 1)
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverTruncateMethodWithTheListOfTableNames() {
	suite.mockDrivers[0].truncateDefaultCall.Unset()
	suite.mockDrivers[1].truncateDefaultCall.Unset()
	suite.mockDrivers[0].On("Truncate", []string{"first table", "another table"}).Return(nil)
	suite.mockDrivers[1].On("Truncate", []string{"first table", "another table"}).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["first table"] = []interface{}{"some object"}
	fixtures["another table"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"first table", "another table"}, fixtures)

	suite.mockDrivers[0].AssertExpectations(suite.T())
	suite.mockDrivers[1].AssertExpectations(suite.T())
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverInsertFixturesMethodOnceForEachTable() {
	fixtures := make(map[string][]interface{})
	fixtures["first table"] = []interface{}{"some object"}
	fixtures["another table"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"first table", "another table"}, fixtures)

	suite.mockDrivers[0].AssertNumberOfCalls(suite.T(), "InsertFixtures", 2)
	suite.mockDrivers[1].AssertNumberOfCalls(suite.T(), "InsertFixtures", 2)
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverInsertFixturesMethodWithEachTableFixtures() {
	suite.mockDrivers[0].insertFixturesDefaultCall.Unset()
	suite.mockDrivers[1].insertFixturesDefaultCall.Unset()

	suite.mockDrivers[0].On("InsertFixtures", "table1", []interface{}{"some object"}).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", "table1", []interface{}{"some object"}).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", "table5", []interface{}{"another object"}).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", "table5", []interface{}{"another object"}).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["table1"] = []interface{}{"some object"}
	fixtures["table5"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"table1", "table5"}, fixtures)

	suite.mockDrivers[0].AssertExpectations(suite.T())
	suite.mockDrivers[1].AssertExpectations(suite.T())
}

func TestDbfixturesSuite(t *testing.T) {
	suite.Run(t, new(dbfixturesTestSuite))
}

/*
mockDriver is a mock of the Driver struct to be used in testing.
*/
type mockDriver struct {
	mock.Mock
	truncateDefaultCall       *mock.Call
	insertFixturesDefaultCall *mock.Call
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
