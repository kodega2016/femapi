package main

import (
	"github.com/kodega2016/femapi/internal/app"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Panicln("we are running our application...")
}
