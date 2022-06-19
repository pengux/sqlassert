package sqlassert

import (
	"database/sql"
)

const (
	errTableNotExists      = "table '%s' does not exist"
	errTableExists         = "table '%s' exists"
	errColumnNotExists     = "column '%s' does not exist in table '%s'"
	errColumnExists        = "column '%s' exists in table '%s'"
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

type nilTestingT struct{}

func (n nilTestingT) Errorf(format string, args ...interface{}) {}

var nilT = new(nilTestingT)

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

func tableExists(
	t testingT,
	db *sql.DB,
	query, table string,
) bool {
	exists := queryExists(db, query, table)
	if !exists {
		t.Errorf(errTableNotExists, table)
	}

	return exists
}

func tableNotExists(
	t testingT,
	db *sql.DB,
	query, table string,
) bool {
	exists := queryExists(db, query, table)
	if exists {
		t.Errorf(errTableExists, table)
	}

	return !exists
}

func columnExists(
	t testingT,
	db *sql.DB,
	query, table, column string,
) bool {
	exists := queryExists(db, query, table, column)
	if !exists {
		t.Errorf(errColumnNotExists, column, table)
	}

	return exists
}

func columnNotExists(
	t testingT,
	db *sql.DB,
	query, table, column string,
) bool {
	exists := queryExists(db, query, table, column)
	if exists {
		t.Errorf(errColumnExists, column, table)
	}

	return !exists
}

func constraintExists(
	t testingT,
	db *sql.DB,
	query, table, constraint string,
) bool {
	exists := queryExists(db, query, table, constraint)
	if !exists {
		t.Errorf(errConstraintNotExists, constraint, table)
	}

	return exists
}

func constraintNotExists(
	t testingT,
	db *sql.DB,
	query, table, constraint string,
) bool {
	exists := queryExists(db, query, table, constraint)
	if exists {
		t.Errorf(errConstraintExists, constraint, table)
	}

	return !exists
}

type rowExistsQueryBuilder func(table string, colVals map[string]interface{}) (string, []interface{})

func rowExists(
	t testingT,
	db *sql.DB,
	table string,
	colVals map[string]interface{},
	builder rowExistsQueryBuilder,
) bool {
	query, args := builder(table, colVals)

	exists := queryExists(db, query, args...)
	if !exists {
		t.Errorf(errRowNotExists, colVals, table)
	}

	return exists
}

func rowNotExists(
	t testingT,
	db *sql.DB,
	table string,
	colVals map[string]interface{},
	builder rowExistsQueryBuilder,
) bool {
	query, args := builder(table, colVals)

	exists := queryExists(db, query, args...)
	if exists {
		t.Errorf(errRowExists, colVals, table)
	}

	return !exists
}

func indexExists(
	t testingT,
	db *sql.DB,
	query string, table, index string,
) bool {
	exists := queryExists(db, query, table, index)
	if !exists {
		t.Errorf(errIndexNotExists, index, table)
	}

	return exists
}

func indexNotExists(
	t testingT,
	db *sql.DB,
	query, table, index string,
) bool {
	exists := queryExists(db, query, table, index)
	if exists {
		t.Errorf(errIndexExists, index, table)
	}

	return !exists
}
