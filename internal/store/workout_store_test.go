package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}
	// run the migration for the test db
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error:%v", err)
	}

	// clearing databases
	_, err = db.Exec("TRUNCATE workouts,workouts_entries CASCADE")
	if err != nil {
		t.Fatalf("truncate database:%v", err)
	}

	return db
}
