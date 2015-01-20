package booking

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Application is the primary exported object
type Application struct {
	registry *Registry
}

// Service Locator for dependencies
type Registry struct {
	db *sql.DB
}

func (r *Registry) DB() *sql.DB {
	return r.db
}

func (r *Registry) Close() error {
	return r.db.Close()
}

// Create an Application with an attached SQL dataSource
// In this case, a path to a sqlite3 database
func NewApplication(dataSourceName string) *Application {
	app := &Application{}
	registry := &Registry{}
	app.registry = registry

	var err error
	registry.db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		panic(err)
	}

	return app
}

// NewServer returns an http.Server for the guest facing booking web UI
func (app *Application) NewServer(addr string) http.Server {
	s := http.Server{}
	s.Addr = addr
	s.Handler = newBookingController(app.registry)
	return s
}

func (app *Application) Close() error {
	err := app.registry.Close()
	if err != nil {
		log.Print(applicationError{err})
		return err
	}

	return nil
}

type applicationError struct {
	err error
}

func (e applicationError) Error() string {
	return "[Application Error] " + e.err.Error()
}
