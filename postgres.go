package sqlassert

import (
	"database/sql"
	"fmt"
	"strings"
)

// PostgresAsserter contains assertion helpers for Postgresql databases
type PostgresAsserter struct {
	db *sql.DB
}

func NewPostgresAsserter(db *sql.DB) *PostgresAsserter {
	return &PostgresAsserter{db}
}

func (pa *PostgresAsserter) TableExists(t testingT, table string) bool {
	exists := false
	err := pa.db.QueryRow(postgresTableExistsQuery, table).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if !exists {
		t.Errorf("table %s does not exist", table)
	}

	return exists
}

func (pa *PostgresAsserter) TableNotExists(t testingT, table string) bool {
	exists := pa.TableExists(nilT, table)

	if exists {
		t.Errorf("table %s exists", table)
	}

	return !exists
}

func (pa *PostgresAsserter) ColumnExists(t testingT, table, column string) bool {
	exists := false
	err := pa.db.QueryRow(postgresColumnExistsQuery, table, column).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if !exists {
		t.Errorf("column %s does not exist in table %s", column, table)
	}

	return exists
}

func (pa *PostgresAsserter) ColumnNotExists(t testingT, table, column string) bool {
	exists := pa.ColumnExists(nilT, table, column)

	if exists {
		t.Errorf("column %s exists in table %s", column, table)
	}

	return !exists
}

func (pa *PostgresAsserter) ConstraintExists(t testingT, table, constraint string) bool {
	exists := false
	err := pa.db.QueryRow(postgresConstraintExistsQuery, table, constraint).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if !exists {
		t.Errorf("constraint %s does not exist in table %s", constraint, table)
	}

	return exists
}

func (pa *PostgresAsserter) ConstraintNotExists(t testingT, table, constraint string) bool {
	exists := pa.ConstraintExists(nilT, table, constraint)

	if exists {
		t.Errorf("constraint %s exists in table %s", constraint, table)
	}

	return !exists
}

func (pa *PostgresAsserter) RowExists(t testingT, table string, colVals map[string]interface{}) bool {
	var wheres []string
	var vals []interface{}

	for col, val := range colVals {
		wheres = append(wheres, fmt.Sprintf(`"%s" = $%d`, col, len(wheres)+1))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf(postgresRowExistsQuery, table, strings.Join(wheres, " AND "))

	exists := false
	err := pa.db.QueryRow(sql, vals...).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if !exists {
		t.Errorf("row with criteria %v does not exist in table %s", colVals, table)
	}

	return exists
}

func (pa *PostgresAsserter) RowNotExists(t testingT, table string, colVals map[string]interface{}) bool {
	exists := pa.RowExists(nilT, table, colVals)

	if exists {
		t.Errorf("row with criteria %v exists in table %s", colVals, table)
	}

	return !exists
}

// IndexExists returns true if the index exists for the table
func (pa *PostgresAsserter) IndexExists(t testingT, table, index string) bool {
	exists := false
	err := pa.db.QueryRow(postgresIndexExistsQuery, table, index).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if !exists {
		t.Errorf("index %s does not exist in table %s", index, table)
	}

	return exists
}

func (pa *PostgresAsserter) IndexNotExists(t testingT, table, index string) bool {
	exists := pa.IndexExists(nilT, table, index)

	if exists {
		t.Errorf("index %s exists in table %s", index, table)
	}

	return !exists
}
