package dbfixtures_test

import (
	"errors"

	"github.com/PedroHenriques/go-dbfixtures/dbfixtures"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type dbfixturesUnitTestSuite struct {
	suite.Suite
	mockDrivers []*mockDriver
}

func (suite *dbfixturesUnitTestSuite) SetupTest() {
	driver1 := &mockDriver{}
	driver2 := &mockDriver{}

	suite.mockDrivers = []*mockDriver{driver1, driver2}
}

func (suite *dbfixturesUnitTestSuite) TestNewItShouldReturnAnInstanceOfTheInterfaceIDbfixtures() {
	res := dbfixtures.New()

	require.Implements(suite.T(), (*dbfixtures.IDbfixtures)(nil), res)
}

func (suite *dbfixturesUnitTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverTruncateMethodOnce() {
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

func (suite *dbfixturesUnitTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverTruncateMethodWithTheListOfTableNames() {
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

func (suite *dbfixturesUnitTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverInsertFixturesMethodOnceForEachTable() {
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

func (suite *dbfixturesUnitTestSuite) TestDbfixturesInsertFixturesItShouldCallEachDriverInsertFixturesMethodWithEachTableFixtures() {
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

func (suite *dbfixturesUnitTestSuite) TestDbfixturesInsertFixturesItShouldReturnAnErrorIfACallToTheDriverTruncateReturnAnError() {
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

func (suite *dbfixturesUnitTestSuite) TestDbfixturesInsertFixturesItShouldReturnAnErrorIfACallToTheDriverInsertFixturesReturnAnError() {
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

func (suite *dbfixturesUnitTestSuite) TestDbfixturesCloseDriversItShouldCallEachDriverCloseMethodOnce() {
	suite.mockDrivers[0].On("Close").Return(nil)
	suite.mockDrivers[1].On("Close").Return(nil)

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	fixturesHandler.CloseDrivers()

	suite.mockDrivers[0].AssertNumberOfCalls(suite.T(), "Close", 1)
	suite.mockDrivers[1].AssertNumberOfCalls(suite.T(), "Close", 1)
}

func (suite *dbfixturesUnitTestSuite) TestDbfixturesCloseDriversItShouldReturnAnErrorIfTheCallToADriverCloseReturnsAnError() {
	suite.mockDrivers[0].On("Close").Return(errors.New("error from the test close()"))
	suite.mockDrivers[1].On("Close").Return(nil)

	fixturesHandler := dbfixtures.New(suite.mockDrivers[0], suite.mockDrivers[1])
	res := fixturesHandler.CloseDrivers()

	require.EqualError(suite.T(), res, "error from the test close()")
}
