package booking

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Application is the primary exported object and doubles as a Registry
type Application struct {
	db *sql.DB
}

// Service Locator for dependencies
// TODO perhaps make more abstract like Registry.Get()/Set()
type Registry interface {
	DB() *sql.DB
}

// Create an Application with an attached SQL dataSource
// In this case, a path to a sqlite3 database
func NewApplication(dataSourceName string) *Application {
	app := &Application{}

	var err error
	app.db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		panic(err)
	}

	return app
}

// NewServer returns an http.Server for the guest facing booking web UI
func (app *Application) NewServer(addr string) http.Server {
	s := http.Server{}
	s.Addr = addr
	s.Handler = newBookingController(app)
	return s
}

func (app *Application) DB() *sql.DB {
	return app.db
}

func (app *Application) Close() error {
	err := app.db.Close()
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
