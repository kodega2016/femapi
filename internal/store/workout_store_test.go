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

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name      string
		workout   *Workout
		wantError bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:             "push day",
				Description:       "upper body day",
				DurationInMinutes: 60,
				CaloriesBurned:    200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						ExerciseSets: 2,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(135.5),
						Notes:        "warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantError: false,
		},
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float32) *float32 {
	return &f
}
