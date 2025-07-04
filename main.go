package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kodega2016/femapi/internal/app"
	"github.com/kodega2016/femapi/internal/routes"
	"github.com/newrelic/go-agent/v3/newrelic"

	_ "github.com/newrelic/go-agent/v3/newrelic"
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

	// close the database
	defer app.DB.Close()

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

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to load env:", err)
	}
	_, err = newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("NEW_RELIC_APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_KEY")),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		log.Fatal("failed to initialize new relic:", err)
	}
}
