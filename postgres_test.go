package sqlassert_test

import (
	"fmt"
	"testing"

	"github.com/pengux/sqlassert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresAsserter(t *testing.T) {
	pgasserter := sqlassert.NewPostgresAsserter(getPostgresDB())

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

func TestPostgresAsserterError(t *testing.T) {
	pgasserter := sqlassert.NewPostgresAsserter(getPostgresDB())
	mockT := new(mockT)

	table := "sqlassert_test"
	column := "sku"
	constraint := "sqlassert_test_pkey"
	row := map[string]interface{}{"sku": "sku1", "name": "name1"}
	index := "sqlassert_test_name_idx"
	nonExisting := "non_existing"

	pgasserter.TableExists(mockT, nonExisting)
	mockT.expectLastError(t, "table '"+nonExisting+"' does not exist")

	pgasserter.ColumnExists(mockT, table, nonExisting)
	mockT.expectLastError(t, "column '"+nonExisting+"' does not exist in table '"+table+"'")

	pgasserter.ConstraintExists(mockT, table, nonExisting)
	mockT.expectLastError(t, "constraint '"+nonExisting+"' does not exist in table '"+table+"'")

	pgasserter.RowExists(mockT, table, map[string]interface{}{"sku": nonExisting})
	mockT.expectLastError(t, "row with criteria map[sku:"+nonExisting+"] does not exist in table '"+table+"'")

	pgasserter.IndexExists(mockT, table, nonExisting)
	mockT.expectLastError(t, "index '"+nonExisting+"' does not exist in table '"+table+"'")

	pgasserter.TableNotExists(mockT, table)
	mockT.expectLastError(t, "table '"+table+"' exists")

	pgasserter.ColumnNotExists(mockT, table, column)
	mockT.expectLastError(t, "column '"+column+"' exists in table '"+table+"'")

	pgasserter.ConstraintNotExists(mockT, table, constraint)
	mockT.expectLastError(t, "constraint '"+constraint+"' exists in table '"+table+"'")

	pgasserter.RowNotExists(mockT, table, row)
	mockT.expectLastError(t, "row with criteria map[name:name1 sku:sku1] exists in table '"+table+"'")

	pgasserter.IndexNotExists(mockT, table, index)
	mockT.expectLastError(t, "index '"+index+"' exists in table '"+table+"'")
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
