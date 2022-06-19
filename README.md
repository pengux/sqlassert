![CI](https://github.com/pengux/sqlassert/actions/workflows/ci.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/pengux/sqlassert.svg)](https://pkg.go.dev/github.com/pengux/sqlassert)
[![Go Report Card](https://goreportcard.com/badge/github.com/pengux/sqlassert)](https://goreportcard.com/report/github.com/pengux/sqlassert)

# sqlassert
Go package with assertion helpers for SQL databases

Use `sqlassert` in your integration tests or write assertions for schema changes and migraton scripts.

## Usage

```
go get github.com/pengux/sqlassert
```

Example:

```
package migration_test

import (
	"testing"

	"github.com/pengux/sqlassert"
)

func TestMigration(t *testing.T) {
	...
	// db is *sql.DB
	pgassert := sqlassert.NewPostgresAsserter(db)

	table := "table_name"
	column := "column_name"
	index := "index_name"
	constraint := "constraint_name"
	id := "123"

	pgassert.TableExists(t, table)
	pgassert.ColumnExists(t, table, column)
	pgassert.ConstraintExists(t, table, column, constraint)
	pgassert.RowExists(t, table, map[string]interface{}{
		"id": id,
	})
	pgassert.IndexExists(t, table, index)
}
```

## Supported databases

- [x] Postgresql
- [x] Mysql
- [ ] SQLite
- [ ] Microsoft SQL Server
- [ ] Oracle
