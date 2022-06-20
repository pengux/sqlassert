package sqlassert

import (
	"database/sql"
	"fmt"
	"strings"
)

// MysqlAsserter contains assertion helpers for mysqlql databases
type MysqlAsserter struct {
	db                    *sql.DB
	dbName                string
	rowExistsQueryBuilder rowExistsQueryBuilder
}

func NewMysqlAsserter(db *sql.DB) *MysqlAsserter {
	var dbName string
	err := db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		panic(err)
	}

	return &MysqlAsserter{
		db:     db,
		dbName: dbName,
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
	exists := queryExists(pa.db, mysqlTableExistsQuery, pa.dbName, table)
	if !exists {
		t.Errorf(errTableNotExists, table)
	}

	return exists
}

func (pa *MysqlAsserter) TableNotExists(t testingT, table string) bool {
	exists := queryExists(pa.db, mysqlTableExistsQuery, pa.dbName, table)
	if exists {
		t.Errorf(errTableExists, table)
	}

	return !exists
}

func (pa *MysqlAsserter) ColumnExists(t testingT, table, column string) bool {
	exists := queryExists(pa.db, mysqlColumnExistsQuery, pa.dbName, table, column)
	if !exists {
		t.Errorf(errColumnNotExists, column, table)
	}

	return exists
}

func (pa *MysqlAsserter) ColumnNotExists(t testingT, table, column string) bool {
	exists := queryExists(pa.db, mysqlColumnExistsQuery, pa.dbName, table, column)
	if exists {
		t.Errorf(errColumnExists, column, table)
	}

	return !exists
}

func (pa *MysqlAsserter) ColumnOfType(t testingT, table, column, dataType string) bool {
	exists := queryExists(pa.db, mysqlColumnOfTypeQuery, pa.dbName, table, column, dataType)
	if !exists {
		t.Errorf(errColumnNotOfType, column, table, dataType)
	}

	return exists
}

func (pa *MysqlAsserter) ConstraintExists(t testingT, table, constraint string) bool {
	exists := queryExists(pa.db, mysqlConstraintExistsQuery, pa.dbName, table, constraint)
	if !exists {
		t.Errorf(errConstraintNotExists, constraint, table)
	}

	return exists
}

func (pa *MysqlAsserter) ConstraintNotExists(t testingT, table, constraint string) bool {
	exists := queryExists(pa.db, mysqlConstraintExistsQuery, pa.dbName, table, constraint)
	if exists {
		t.Errorf(errConstraintExists, constraint, table)
	}

	return !exists
}

func (pa *MysqlAsserter) RowExists(t testingT, table string, colVals map[string]interface{}) bool {
	query, args := pa.rowExistsQueryBuilder(table, colVals)

	exists := queryExists(pa.db, query, args...)
	if !exists {
		t.Errorf(errRowNotExists, colVals, table)
	}

	return exists
}

func (pa *MysqlAsserter) RowNotExists(t testingT, table string, colVals map[string]interface{}) bool {
	query, args := pa.rowExistsQueryBuilder(table, colVals)

	exists := queryExists(pa.db, query, args...)
	if exists {
		t.Errorf(errRowExists, colVals, table)
	}

	return !exists
}

func (pa *MysqlAsserter) IndexExists(t testingT, table, index string) bool {
	exists := queryExists(pa.db, mysqlIndexExistsQuery, pa.dbName, table, index)
	if !exists {
		t.Errorf(errIndexNotExists, index, table)
	}

	return exists
}

func (pa *MysqlAsserter) IndexNotExists(t testingT, table, index string) bool {
	exists := queryExists(pa.db, mysqlIndexExistsQuery, pa.dbName, table, index)
	if exists {
		t.Errorf(errIndexExists, index, table)
	}

	return !exists
}

const (
	mysqlTableExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ?
)
`
	mysqlColumnExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.columns
	WHERE
		table_schema = ?
		AND table_name = ?
		AND column_name = ?
);
`
	mysqlColumnOfTypeQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.columns
	WHERE
		table_schema = ?
		AND table_name = ?
		AND column_name = ?
		AND data_type = ?
);
`
	mysqlConstraintExistsQuery = `
SELECT EXISTS(
	SELECT 1 FROM information_schema.table_constraints
	WHERE
		table_schema = ?
		AND table_name = ?
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
		table_schema = ?
		AND table_name = ?
		AND index_name = ?
);
`
)
