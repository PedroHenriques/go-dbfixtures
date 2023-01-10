package dbfixtures_test

import (
	"errors"
	"testing"

	"github.com/PedroHenriques/go-dbfixtures/dbfixtures"
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

	suite.mockDrivers = []*mockDriver{driver1, driver2}
}

func (suite *dbfixturesTestSuite) TestNewItShouldReturnAnInstanceOfTheTypeDbfixtures() {
	res := dbfixtures.New()

	require.Implements(suite.T(), (*dbfixtures.IDbfixtures)(nil), res)
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverTruncateMethodOnce() {
	suite.mockDrivers[0].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[1].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["table1"] = []interface{}{"some object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"table1"}, fixtures)

	suite.mockDrivers[0].AssertNumberOfCalls(suite.T(), "Truncate", 1)
	suite.mockDrivers[1].AssertNumberOfCalls(suite.T(), "Truncate", 1)
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverTruncateMethodWithTheListOfTableNames() {
	suite.mockDrivers[0].On("Truncate", []string{"first table", "another table"}).Return(nil)
	suite.mockDrivers[1].On("Truncate", []string{"first table", "another table"}).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["first table"] = []interface{}{"some object"}
	fixtures["another table"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"first table", "another table"}, fixtures)

	suite.mockDrivers[0].AssertExpectations(suite.T())
	suite.mockDrivers[1].AssertExpectations(suite.T())
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverInsertFixturesMethodOnceForEachTable() {
	suite.mockDrivers[0].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[1].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["first table"] = []interface{}{"some object"}
	fixtures["another table"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"first table", "another table"}, fixtures)

	suite.mockDrivers[0].AssertNumberOfCalls(suite.T(), "InsertFixtures", 2)
	suite.mockDrivers[1].AssertNumberOfCalls(suite.T(), "InsertFixtures", 2)
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverInsertFixturesMethodWithEachTableFixtures() {
	suite.mockDrivers[0].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[1].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", "table1", []interface{}{"some object"}).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", "table5", []interface{}{"another object"}).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", "table1", []interface{}{"some object"}).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", "table5", []interface{}{"another object"}).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["table1"] = []interface{}{"some object"}
	fixtures["table5"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.InsertFixtures([]string{"table1", "table5"}, fixtures)

	suite.mockDrivers[0].AssertExpectations(suite.T())
	suite.mockDrivers[1].AssertExpectations(suite.T())
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldReturnAnErrorIfACallToTheDriverTruncateReturnAnError() {
	suite.mockDrivers[0].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[1].On("Truncate", []string{"first table", "another table"}).Return(errors.New("Error from the unit test"))
	suite.mockDrivers[0].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)
	suite.mockDrivers[1].On("InsertFixtures", mock.Anything, mock.Anything).Return(nil)

	fixtures := make(map[string][]interface{})
	fixtures["first table"] = []interface{}{"some object"}
	fixtures["another table"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	res := fixturesHandler.InsertFixtures([]string{"first table", "another table"}, fixtures)

	require.EqualError(suite.T(), res, "Error from the unit test")
}

func (suite *dbfixturesTestSuite) TestDbfixturesInsertFixturesItShouldReturnAnErrorIfACallToTheDriverInsertFixturesReturnAnError() {
	suite.mockDrivers[0].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[1].On("Truncate", mock.Anything).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", "table1", []interface{}{"some object"}).Return(nil)
	suite.mockDrivers[0].On("InsertFixtures", "table5", []interface{}{"another object"}).Return(errors.New("Another error from the unit test."))

	fixtures := make(map[string][]interface{})
	fixtures["table1"] = []interface{}{"some object"}
	fixtures["table5"] = []interface{}{"another object"}

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	res := fixturesHandler.InsertFixtures([]string{"table1", "table5"}, fixtures)

	require.EqualError(suite.T(), res, "Another error from the unit test.")
}

func (suite *dbfixturesTestSuite) TestDbfixturesCloseDriversItShouldCallEachDriverCloseMethodOnce() {
	suite.mockDrivers[0].On("Close").Return(nil)
	suite.mockDrivers[1].On("Close").Return(nil)

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.CloseDrivers()

	suite.mockDrivers[0].AssertNumberOfCalls(suite.T(), "Close", 1)
	suite.mockDrivers[1].AssertNumberOfCalls(suite.T(), "Close", 1)
}

func (suite *dbfixturesTestSuite) TestDbfixturesCloseDriversItShouldReturnAnErrorIfTheCallToADriverCloseReturnsAnError() {
	suite.mockDrivers[0].On("Close").Return(errors.New("error from the test close()"))
	suite.mockDrivers[1].On("Close").Return(nil)

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	res := fixturesHandler.CloseDrivers()

	require.EqualError(suite.T(), res, "error from the test close()")
}

func TestDbfixturesSuite(t *testing.T) {
	suite.Run(t, new(dbfixturesTestSuite))
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