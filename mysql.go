package sqlassert

import (
	"database/sql"
	"fmt"
	"strings"
)

// MysqlAsserter contains assertion helpers for mysqlql databases
type MysqlAsserter struct {
	db                    *sql.DB
	rowExistsQueryBuilder rowExistsQueryBuilder
}

func NewMysqlAsserter(db *sql.DB) *MysqlAsserter {
	return &MysqlAsserter{
		db: db,
		rowExistsQueryBuilder: func(
			table string,
			colVals map[string]interface{},
		) (string, []interface{}) {
			var wheres []string
			var vals []interface{}

			for col, val := range colVals {
				wheres = append(wheres, fmt.Sprintf("`%s` = ?", col))
				vals = append(vals, val)
			}
			sql := fmt.Sprintf(mysqlRowExistsQuery, table, strings.Join(wheres, " AND "))

			return sql, vals
		},
	}
}

func (pa *MysqlAsserter) TableExists(t testingT, table string) bool {
	return tableExists(t, pa.db, mysqlTableExistsQuery, table)
}

func (pa *MysqlAsserter) TableNotExists(t testingT, table string) bool {
	return tableNotExists(t, pa.db, mysqlTableExistsQuery, table)
}

func (pa *MysqlAsserter) ColumnExists(t testingT, table, column string) bool {
	return columnExists(t, pa.db, mysqlColumnExistsQuery, table, column)
}

func (pa *MysqlAsserter) ColumnNotExists(t testingT, table, column string) bool {
	return columnNotExists(t, pa.db, mysqlColumnExistsQuery, table, column)
}

func (pa *MysqlAsserter) ConstraintExists(t testingT, table, constraint string) bool {
	return constraintExists(t, pa.db, mysqlConstraintExistsQuery, table, constraint)
}

func (pa *MysqlAsserter) ConstraintNotExists(t testingT, table, constraint string) bool {
	return constraintNotExists(t, pa.db, mysqlConstraintExistsQuery, table, constraint)
}

func (pa *MysqlAsserter) RowExists(t testingT, table string, colVals map[string]interface{}) bool {
	return rowExists(t, pa.db, table, colVals, pa.rowExistsQueryBuilder)
}

func (pa *MysqlAsserter) RowNotExists(t testingT, table string, colVals map[string]interface{}) bool {
	return rowNotExists(t, pa.db, table, colVals, pa.rowExistsQueryBuilder)
}

func (pa *MysqlAsserter) IndexExists(t testingT, table, index string) bool {
	return indexExists(t, pa.db, mysqlIndexExistsQuery, table, index)
}

func (pa *MysqlAsserter) IndexNotExists(t testingT, table, index string) bool {
	return indexNotExists(t, pa.db, mysqlIndexExistsQuery, table, index)
}

const (
	mysqlTableExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM information_schema.tables WHERE table_name = ?
)
`
	mysqlColumnExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.columns
	WHERE
		table_name = ?
		AND column_name = ?
);
`
	mysqlConstraintExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.table_constraints
	WHERE
		table_name = ?
		AND constraint_name = ?
);
`

	mysqlRowExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM ` + "`%s`" + ` WHERE %s
);
`

	mysqlIndexExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.statistics
	WHERE
		table_name = ?
		AND index_name = ?
);
`
)
