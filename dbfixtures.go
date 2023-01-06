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
		driver.Truncate(tableNames)

		for _, tableName := range tableNames {
			driver.InsertFixtures(tableName, fixtures[tableName])
		}
	}

	return nil
}
