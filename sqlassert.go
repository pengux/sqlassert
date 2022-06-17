package sqlassert

import (
	"database/sql"
)

type testingT interface {
	Errorf(format string, args ...interface{})
}

type nilTestingT struct{}

func (n nilTestingT) Errorf(format string, args ...interface{}) {}

var nilT = new(nilTestingT)

func queryExists(
	t testingT,
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
	exists := queryExists(t, db, query, table)
	if !exists {
		t.Errorf("table %s does not exist", table)
	}

	return exists
}

func tableNotExists(
	t testingT,
	db *sql.DB,
	query, table string,
) bool {
	exists := queryExists(nilT, db, query, table)
	if exists {
		t.Errorf("table %s exists", table)
	}

	return !exists
}

func columnExists(
	t testingT,
	db *sql.DB,
	query, table, column string,
) bool {
	exists := queryExists(t, db, query, table, column)
	if !exists {
		t.Errorf("column %s does not exist in table %s", column, table)
	}

	return exists
}

func columnNotExists(
	t testingT,
	db *sql.DB,
	query, table, column string,
) bool {
	exists := queryExists(nilT, db, query, table, column)
	if exists {
		t.Errorf("column %s exists in table %s", column, table)
	}

	return !exists
}

func constraintExists(
	t testingT,
	db *sql.DB,
	query, table, constraint string,
) bool {
	exists := queryExists(t, db, query, table, constraint)
	if !exists {
		t.Errorf("constraint %s does not exist in table %s", constraint, table)
	}

	return exists
}

func constraintNotExists(
	t testingT,
	db *sql.DB,
	query, table, constraint string,
) bool {
	exists := queryExists(nilT, db, query, table, constraint)
	if exists {
		t.Errorf("constraint %s exists in table %s", constraint, table)
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

	exists := queryExists(nilT, db, query, args...)
	if !exists {
		t.Errorf("row with criteria %v does not exist in table %s", colVals, table)
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

	exists := queryExists(nilT, db, query, args...)
	if exists {
		t.Errorf("row with criteria %v exists in table %s", colVals, table)
	}

	return !exists
}

func indexExists(
	t testingT,
	db *sql.DB,
	query, table, index string,
) bool {
	exists := queryExists(nilT, db, query, table, index)
	if !exists {
		t.Errorf("index %s does not exist in table %s", index, table)
	}

	return exists
}

func indexNotExists(
	t testingT,
	db *sql.DB,
	query, table, index string,
) bool {
	exists := queryExists(nilT, db, query, table, index)
	if exists {
		t.Errorf("index %s exists in table %s", index, table)
	}

	return !exists
}
