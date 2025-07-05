package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	_, err = db.Exec("TRUNCATE workouts,workout_entries CASCADE")
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
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:             "full body",
				Description:       "complete workout",
				DurationInMinutes: 90,
				CaloriesBurned:    500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						ExerciseSets: 3,
						Reps:         IntPtr(60),
						Notes:        "keep form",
						OrderIndex:   1,
					}, {
						ExerciseName: "sqauts",
						ExerciseSets: 4,
						Reps:         IntPtr(12),
						Notes:        "full depth",
						OrderIndex:   2,
					},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationInMinutes, createdWorkout.DurationInMinutes)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.Title, retrieved.Title)
			assert.Equal(t, createdWorkout.Description, retrieved.Description)
			assert.Equal(t, createdWorkout.DurationInMinutes, retrieved.DurationInMinutes)
			assert.Equal(t, len(createdWorkout.Entries), len(retrieved.Entries))
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float32) *float32 {
	return &f
}
