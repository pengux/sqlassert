package sqlassert_test

import (
	"fmt"
	"testing"

	"github.com/pengux/sqlassert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMysqlAsserter(t *testing.T) {
	mysqlAsserter := sqlassert.NewMysqlAsserter(getMysqlDB())

	table := "sqlassert_test"
	column := "sku"
	dataType := "varchar"
	constraint := "PRIMARY"
	row := map[string]interface{}{"sku": "sku1", "name": "name1"}
	index := "sqlassert_test_name_idx"
	nonExisting := "non_existing"

	mysqlAsserter.TableExists(t, table)
	mysqlAsserter.ColumnExists(t, table, column)
	mysqlAsserter.ColumnOfType(t, table, column, dataType)
	mysqlAsserter.ConstraintExists(t, table, constraint)
	mysqlAsserter.RowExists(t, table, row)
	mysqlAsserter.IndexExists(t, table, index)

	mysqlAsserter.TableNotExists(t, nonExisting)
	mysqlAsserter.ColumnNotExists(t, table, nonExisting)
	mysqlAsserter.ConstraintNotExists(t, table, nonExisting)
	mysqlAsserter.RowNotExists(t, table, map[string]interface{}{"sku": nonExisting})
	mysqlAsserter.IndexNotExists(t, table, nonExisting)
}

func TestMysqlAsserterError(t *testing.T) {
	mysqlAsserter := sqlassert.NewMysqlAsserter(getMysqlDB())
	mockT := new(mockT)

	table := "sqlassert_test"
	column := "sku"
	constraint := "PRIMARY"
	row := map[string]interface{}{"sku": "sku1", "name": "name1"}
	index := "sqlassert_test_name_idx"
	nonExisting := "non_existing"

	mysqlAsserter.TableExists(mockT, nonExisting)
	mockT.expectLastError(t, "table '"+nonExisting+"' does not exist")

	mysqlAsserter.ColumnExists(mockT, table, nonExisting)
	mockT.expectLastError(t, "column '"+nonExisting+"' does not exist in table '"+table+"'")

	mysqlAsserter.ColumnOfType(mockT, table, column, nonExisting)
	mockT.expectLastError(t, "column '"+column+"' in table '"+table+"' is not of type '"+nonExisting+"'")

	mysqlAsserter.ConstraintExists(mockT, table, nonExisting)
	mockT.expectLastError(t, "constraint '"+nonExisting+"' does not exist in table '"+table+"'")

	mysqlAsserter.RowExists(mockT, table, map[string]interface{}{"sku": nonExisting})
	mockT.expectLastError(t, "row with criteria map[sku:"+nonExisting+"] does not exist in table '"+table+"'")

	mysqlAsserter.IndexExists(mockT, table, nonExisting)
	mockT.expectLastError(t, "index '"+nonExisting+"' does not exist in table '"+table+"'")

	mysqlAsserter.TableNotExists(mockT, table)
	mockT.expectLastError(t, "table '"+table+"' exists")

	mysqlAsserter.ColumnNotExists(mockT, table, column)
	mockT.expectLastError(t, "column '"+column+"' exists in table '"+table+"'")

	mysqlAsserter.ConstraintNotExists(mockT, table, constraint)
	mockT.expectLastError(t, "constraint '"+constraint+"' exists in table '"+table+"'")

	mysqlAsserter.RowNotExists(mockT, table, row)
	mockT.expectLastError(t, "row with criteria map[name:name1 sku:sku1] exists in table '"+table+"'")

	mysqlAsserter.IndexNotExists(mockT, table, index)
	mockT.expectLastError(t, "index '"+index+"' exists in table '"+table+"'")
}

func TestMysqlAsserterPanic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT DATABASE()").WillReturnError(fmt.Errorf("error"))
	assertPanic(t, "error", func() { _ = sqlassert.NewMysqlAsserter(db) })

	mock.ExpectQuery("SELECT DATABASE()").WillReturnRows(sqlmock.NewRows([]string{"database()"}).AddRow("sqlassert"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))
	mock.ExpectQuery("SELECT EXISTS").WillReturnError(fmt.Errorf("error"))

	mysqlAsserter := sqlassert.NewMysqlAsserter(db)

	table := "sqlassert_test"
	column := "sku"
	dataType := "varchar"
	constraint := "sqlassert_test_pkey"
	row := map[string]interface{}{"sku": "sku1", "name": "name1"}
	index := "sqlassert_test_name_idx"

	assertPanic(t, "error", func() { mysqlAsserter.TableExists(t, table) })
	assertPanic(t, "error", func() { mysqlAsserter.ColumnExists(t, table, column) })
	assertPanic(t, "error", func() { mysqlAsserter.ColumnOfType(t, table, column, dataType) })
	assertPanic(t, "error", func() { mysqlAsserter.ConstraintExists(t, table, constraint) })
	assertPanic(t, "error", func() { mysqlAsserter.RowExists(t, table, row) })
	assertPanic(t, "error", func() { mysqlAsserter.IndexExists(t, table, index) })
}
