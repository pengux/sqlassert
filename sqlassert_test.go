package sqlassert_test

import (
	"database/sql"
	"os"
	"testing"
)

func getDB() *sql.DB {
	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		panic("unable to open database: " + err.Error())
	}

	return db
}

func assertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		r := recover()
		if r == nil || r.(error).Error() != msg {
			t.Errorf("expected panic '%s', got '%v'", msg, r)
		}
	}()
	f()
}
