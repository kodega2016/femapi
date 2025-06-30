// Package api will handle workout related functions
package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHanlder struct{}

func NewWorkoutHandler() *WorkoutHanlder {
	return &WorkoutHanlder{}
}

func (wh *WorkoutHanlder) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Printf("This is the correct workout:%d\n", workoutID)
}

func (wh *WorkoutHanlder) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "workout is creating\n")
}
