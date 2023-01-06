package dbfixtures

type IDbfixtures interface {
	InsertFixtures(tableNames []string, fixtures map[string][]interface{}) error
}

type IDriver interface {
	// clears the specified "tables" from any content
	Truncate(tableNames []string) error

	// inserts the supplied "rows" into the specified "table"
	InsertFixtures(tableName string, fixtures []interface{}) error

	// terminates the connection to the database
	Close() error
}
