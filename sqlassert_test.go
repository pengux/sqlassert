package sqlassert_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type mockT []string

func (mt *mockT) Errorf(format string, args ...interface{}) {
	*mt = append(*mt, fmt.Sprintf(format, args...))
}

func (mt *mockT) lastError() string {
	if len(*mt) == 0 {
		return ""
	}
	return (*mt)[len(*mt)-1]
}

func (mt *mockT) expectLastError(t *testing.T, msg string) {
	t.Helper()
	if mt.lastError() != msg {
		t.Errorf("expected error '%s', got '%s'", msg, mt.lastError())
	}
}

func getPostgresDB() *sql.DB {
	db, err := sql.Open("pgx", os.Getenv("POSTGRES_DB_URL"))
	if err != nil {
		panic("unable to open database: " + err.Error())
	}

	return db
}

func getMysqlDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DB_URL"))
	if err != nil {
		panic("unable to open database: " + err.Error())
	}

	return db
}

func assertPanic(t *testing.T, msg string, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil || r.(error).Error() != msg {
			t.Errorf("expected panic '%s', got '%v'", msg, r)
		}
	}()
	f()
}
