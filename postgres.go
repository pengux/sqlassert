package sqlassert

import (
	"database/sql"
	"fmt"
	"strings"
)

// PostgresAsserter contains assertion helpers for Postgresql databases
type PostgresAsserter struct {
	db                    *sql.DB
	rowExistsQueryBuilder rowExistsQueryBuilder
}

func NewPostgresAsserter(db *sql.DB) *PostgresAsserter {
	return &PostgresAsserter{
		db: db,
		rowExistsQueryBuilder: func(
			table string,
			colVals map[string]interface{},
		) (string, []interface{}) {
			var wheres []string
			var vals []interface{}

			for col, val := range colVals {
				wheres = append(wheres, fmt.Sprintf(`"%s" = $%d`, col, len(wheres)+1))
				vals = append(vals, val)
			}
			sql := fmt.Sprintf(postgresRowExistsQuery, table, strings.Join(wheres, " AND "))

			return sql, vals
		},
	}
}

func (pa *PostgresAsserter) TableExists(t testingT, table string) bool {
	exists := queryExists(pa.db, postgresTableExistsQuery, table)
	if !exists {
		t.Errorf(errTableNotExists, table)
	}

	return exists
}

func (pa *PostgresAsserter) TableNotExists(t testingT, table string) bool {
	exists := queryExists(pa.db, postgresTableExistsQuery, table)
	if exists {
		t.Errorf(errTableExists, table)
	}

	return !exists
}

func (pa *PostgresAsserter) ColumnExists(t testingT, table, column string) bool {
	exists := queryExists(pa.db, postgresColumnExistsQuery, table, column)
	if !exists {
		t.Errorf(errColumnNotExists, column, table)
	}

	return exists
}

func (pa *PostgresAsserter) ColumnNotExists(t testingT, table, column string) bool {
	exists := queryExists(pa.db, postgresColumnExistsQuery, table, column)
	if exists {
		t.Errorf(errColumnExists, column, table)
	}

	return !exists
}

func (pa *PostgresAsserter) ConstraintExists(t testingT, table, constraint string) bool {
	exists := queryExists(pa.db, postgresConstraintExistsQuery, table, constraint)
	if !exists {
		t.Errorf(errConstraintNotExists, constraint, table)
	}

	return exists
}

func (pa *PostgresAsserter) ConstraintNotExists(t testingT, table, constraint string) bool {
	exists := queryExists(pa.db, postgresConstraintExistsQuery, table, constraint)
	if exists {
		t.Errorf(errConstraintExists, constraint, table)
	}

	return !exists
}

func (pa *PostgresAsserter) RowExists(t testingT, table string, colVals map[string]interface{}) bool {
	query, args := pa.rowExistsQueryBuilder(table, colVals)

	exists := queryExists(pa.db, query, args...)
	if !exists {
		t.Errorf(errRowNotExists, colVals, table)
	}

	return exists
}

func (pa *PostgresAsserter) RowNotExists(t testingT, table string, colVals map[string]interface{}) bool {
	query, args := pa.rowExistsQueryBuilder(table, colVals)

	exists := queryExists(pa.db, query, args...)
	if exists {
		t.Errorf(errRowExists, colVals, table)
	}

	return !exists
}

func (pa *PostgresAsserter) IndexExists(t testingT, table, index string) bool {
	exists := queryExists(pa.db, postgresIndexExistsQuery, table, index)
	if !exists {
		t.Errorf(errIndexNotExists, index, table)
	}

	return exists
}

func (pa *PostgresAsserter) IndexNotExists(t testingT, table, index string) bool {
	exists := queryExists(pa.db, postgresIndexExistsQuery, table, index)
	if exists {
		t.Errorf(errIndexExists, index, table)
	}

	return !exists
}

const (
	postgresTableExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM information_schema.tables WHERE table_name = $1
);
`
	postgresColumnExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.columns
	WHERE
		table_name = $1
		AND column_name = $2
);
`
	postgresConstraintExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.constraint_column_usage
	WHERE
		table_name = $1
		AND constraint_name = $2
);
`
	postgresRowExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM "%s" WHERE %s
);
`
	postgresIndexExistsQuery = `
SELECT EXISTS(
	SELECT 1
	FROM
		pg_indexes
	WHERE
		tablename = $1
		AND indexname= $2
);
`
)
