package sqlassert_test

import (
	"fmt"
	"testing"

	"github.com/pengux/sqlassert"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestPostgresAsserter(t *testing.T) {
	pgasserter := sqlassert.NewPostgresAsserter(getDB())

	table := "sqlassert_test"
	column := "sku"
	constraint := "sqlassert_test_pkey"
	row := map[string]interface{}{"sku": "sku1", "name": "name1"}
	index := "sqlassert_test_name_idx"
	nonExisting := "non_existing"

	pgasserter.TableExists(t, table)
	pgasserter.ColumnExists(t, table, column)
	pgasserter.ConstraintExists(t, table, constraint)
	pgasserter.RowExists(t, table, row)
	pgasserter.IndexExists(t, table, index)

	pgasserter.TableNotExists(t, nonExisting)
	pgasserter.ColumnNotExists(t, table, nonExisting)
	pgasserter.ConstraintNotExists(t, table, nonExisting)
	pgasserter.RowNotExists(t, table, map[string]interface{}{"sku": nonExisting})
	pgasserter.IndexNotExists(t, table, nonExisting)
}

func TestPostgresAsserterPanic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))

	pgasserter := sqlassert.NewPostgresAsserter(db)

	table := "sqlassert_test"
	column := "sku"
	constraint := "sqlassert_test_pkey"
	row := map[string]interface{}{"sku": "sku1", "name": "name1"}
	index := "sqlassert_test_name_idx"

	assertPanic(t, "error", func() { pgasserter.TableExists(t, table) })
	assertPanic(t, "error", func() { pgasserter.ColumnExists(t, table, column) })
	assertPanic(t, "error", func() { pgasserter.ConstraintExists(t, table, constraint) })
	assertPanic(t, "error", func() { pgasserter.RowExists(t, table, row) })
	assertPanic(t, "error", func() { pgasserter.IndexExists(t, table, index) })
}
