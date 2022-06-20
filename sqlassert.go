package sqlassert

import (
	"database/sql"
)

const (
	errTableNotExists      = "table '%s' does not exist"
	errTableExists         = "table '%s' exists"
	errColumnNotExists     = "column '%s' does not exist in table '%s'"
	errColumnExists        = "column '%s' exists in table '%s'"
	errColumnNotOfType     = "column '%s' in table '%s' is not of type '%s'"
	errConstraintNotExists = "constraint '%s' does not exist in table '%s'"
	errConstraintExists    = "constraint '%s' exists in table '%s'"
	errRowNotExists        = "row with criteria %v does not exist in table '%s'"
	errRowExists           = "row with criteria %v exists in table '%s'"
	errIndexNotExists      = "index '%s' does not exist in table '%s'"
	errIndexExists         = "index '%s' exists in table '%s'"
)

type testingT interface {
	Errorf(format string, args ...interface{})
}

type rowExistsQueryBuilder func(table string, colVals map[string]interface{}) (string, []interface{})

func queryExists(
	db *sql.DB,
	query string,
	args ...interface{},
) bool {
	exists := false
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
