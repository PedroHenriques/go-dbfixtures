/*
Package dbfixtures provides functionality to seed databases with test fixtures
allowing an easy way to setup DBs before each test to a predictable state.
*/
package dbfixtures

/*
New creates and returns an instance of the DB fixture handler.
*/
func New(drivers ...IDriver) IDbfixtures {
	return &dbfixtures{
		drivers: drivers,
	}
}

type dbfixtures struct {
	drivers []IDriver
}

/*
InsertFixtures will call each registered driver to truncate each "table" and
then insert the relevant data into each "table".
*/
func (handler *dbfixtures) InsertFixtures(
	tableNames []string,
	fixtures map[string][]interface{},
) error {
	for _, driver := range handler.drivers {
		err := driver.Truncate(tableNames)
		if err != nil {
			return err
		}

		for _, tableName := range tableNames {
			err := driver.InsertFixtures(tableName, fixtures[tableName])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/*
CloseDrivers will call each registered driver to perform any necessary cleanup
operations to before the drivers are deleted.
*/
func (handler *dbfixtures) CloseDrivers() error {
	for _, driver := range handler.drivers {
		err := driver.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
