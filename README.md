[![Coverage Status](https://coveralls.io/repos/github/PedroHenriques/go-dbfixtures/badge.svg?branch=main)](https://coveralls.io/github/PedroHenriques/go-dbfixtures?branch=main)
![ci workflow](https://github.com/PedroHenriques/go-dbfixtures/actions/workflows/ci.yml/badge.svg?branch=main)
![cd workflow](https://github.com/PedroHenriques/go-dbfixtures/actions/workflows/cd.yml/badge.svg)

# Fixtures Manager

An abstraction layer for handling database fixtures for automated testing purposes, providing a standardized interface across different database systems.

## Installation

```sh
go get github.com/PedroHenriques/go-dbfixtures
```

## Features

* Test runner agnostic
* No dependencies
* Standardized interface across multiple database systems
* Easily set your database for each test's needs

## Golang versions

- **Package version `1.0.*`** supports Golang `v1.17` or higher.  

## Drivers

This package will use drivers to handle the database operations.
Each driver will be dedicated to 1 databse system (ex: MongoDb, Postgres).  
You can set as many drivers as needed and the fixtures will be sent to each one.

### Driver interface

The drivers are expected to use the following interface

```go
type IDriver interface {
	// clears the specified "tables" of any content
	Truncate(tableNames []string) error

	// inserts the supplied "rows" into the specified "table"
	InsertFixtures(tableName string, fixtures []interface{}) error

	// cleanup and terminate the connection to the database
	Close() error
}
```

### Current official drivers

* [MongoDB](https://github.com/PedroHenriques/go-dbfixtures-mongodb-driver)

## Usage

This package exposes the following function

```go
func New(drivers ...IDriver) IDbfixtures

type IDbfixtures interface {
	InsertFixtures(tableNames []string, fixtures map[string][]interface{}) error
	CloseDrivers() error
}
```

Where

* `InsertFixtures(tableNames []string, fixtures map[string][]interface{}) error`: call this function with the fixtures to be sent to each registered driver.  
**Note:** the fixtures will be inserted in the order they are provided.

* `CloseDrivers() error`: call this function to run any necessary cleanup operations on all registered drivers.

The `IDriver` interface is described in the section **Driver interface** above.

### Example

```go
package driver_test

import (
	"context"

	"github.com/PedroHenriques/go-dbfixtures-mongodb-driver/driver"
	"github.com/PedroHenriques/go-dbfixtures/dbfixtures"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type driverE2eTestSuite struct {
	suite.Suite

	mongoCliente *mongo.Client
	mongoDbName  string
}

func (suite *driverE2eTestSuite) SetupSuite() {
	ConnUrl := "mongodb://testmongo:27017"
	opts := options.Client().ApplyURI(ConnUrl)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	suite.mongoCliente = client

	suite.mongoDbName = "testCol"
}

func (suite *driverE2eTestSuite) TearDownSuite() {
	suite.mongoCliente.Disconnect(context.TODO())
}

func (suite *driverE2eTestSuite) SetupTest() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	err := database.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
}

func (suite *driverE2eTestSuite) TestItShouldWork() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	col1 := database.Collection("col1")

	insRes1, _ := col1.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{Name: "doc 11", Age: 3},
			testDocument{Name: "doc 12", Age: 33},
		},
	)

	require.Equal(suite.T(), 2, len(insRes1.InsertedIDs))

	expectedDocuments := &[]interface{}{
		testDocument{Name: "doc 26", Age: 86},
		testDocument{Name: "doc 27", Age: 87},
	}

	driver := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)

	fixtureHandler := dbfixtures.New(driver)

	err := fixtureHandler.InsertFixtures(
		[]string{"co12"},
		map[string][]interface{}{
			"col1": *expectedDocuments,
		},
	)

	require.Nil(suite.T(), err)

	colDocs, _ := col1.Find(context.TODO(), bson.D{})
	actualDocuments := &[]testDocument{}
	colDocs.All(context.TODO(), actualDocuments)

	require.EqualValues(suite.T(), len(*expectedDocuments), len(*actualDocuments))
	for i, actualDocument := range *actualDocuments {
		require.EqualExportedValues(suite.T(), (*expectedDocuments)[i], actualDocument)
	}
}
```

## How It Works

Each registered driver will be called to:

* clear the "tables" that will be used in the current fixture insertion operation from any content.

* insert the fixtures in the order they were provided.

* terminate the connection to their database.

## Testing This Package

* `cd` into the package's root directory

* Run `sh cli/test.sh -b -gv 1.20`