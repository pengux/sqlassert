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
	return tableExists(t, pa.db, postgresTableExistsQuery, table)
}

func (pa *PostgresAsserter) TableNotExists(t testingT, table string) bool {
	return tableNotExists(t, pa.db, postgresTableExistsQuery, table)
}

func (pa *PostgresAsserter) ColumnExists(t testingT, table, column string) bool {
	return columnExists(t, pa.db, postgresColumnExistsQuery, table, column)
}

func (pa *PostgresAsserter) ColumnNotExists(t testingT, table, column string) bool {
	return columnNotExists(t, pa.db, postgresColumnExistsQuery, table, column)
}

func (pa *PostgresAsserter) ConstraintExists(t testingT, table, constraint string) bool {
	return constraintExists(t, pa.db, postgresConstraintExistsQuery, table, constraint)
}

func (pa *PostgresAsserter) ConstraintNotExists(t testingT, table, constraint string) bool {
	return constraintNotExists(t, pa.db, postgresConstraintExistsQuery, table, constraint)
}

func (pa *PostgresAsserter) RowExists(t testingT, table string, colVals map[string]interface{}) bool {
	return rowExists(t, pa.db, table, colVals, pa.rowExistsQueryBuilder)
}

func (pa *PostgresAsserter) RowNotExists(t testingT, table string, colVals map[string]interface{}) bool {
	return rowNotExists(t, pa.db, table, colVals, pa.rowExistsQueryBuilder)
}

func (pa *PostgresAsserter) IndexExists(t testingT, table, index string) bool {
	return indexExists(t, pa.db, postgresIndexExistsQuery, table, index)
}

func (pa *PostgresAsserter) IndexNotExists(t testingT, table, index string) bool {
	return indexNotExists(t, pa.db, postgresIndexExistsQuery, table, index)
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
		pg_class t,
		pg_class i,
		pg_index ix
	WHERE
		t.oid = ix.indrelid
		AND i.oid = ix.indexrelid
		AND t.relkind = 'r'
		AND t.relname = $1
		AND i.relname = $2
);
`
)
