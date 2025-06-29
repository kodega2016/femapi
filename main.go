package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/kodega2016/femapi/internal/app"
	"github.com/kodega2016/femapi/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "This is the default port on which the server will run")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Printf("we are running our application on port %d\n", port)
	http.HandleFunc("/health", app.HealthCheck)

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
