package dbfixtures

type IDbfixtures interface {
	InsertFixtures(tableNames []string, fixtures map[string][]interface{}) error
	CloseDrivers() error
}

type IDriver interface {
	// clears the specified "tables" of any content
	Truncate(tableNames []string) error

	// inserts the supplied "rows" into the specified "table"
	InsertFixtures(tableName string, fixtures []interface{}) error

	// cleanup and terminate the connection to the database
	Close() error
}
