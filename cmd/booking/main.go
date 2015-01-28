package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cmdrkeene/booking"
	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/facebookgo/inject"
	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

var flagHttp = flag.String("http", ":3000", "http port")

func main() {
	flag.Set("stderrthreshold", "ERROR")
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Domain
	var calendar booking.Calendar
	var formBuilder booking.FormBuilder
	var guestbook booking.Guestbook
	var handler booking.Handler
	var ledger booking.Ledger
	var register booking.Register

	// Dependency Injection
	var g inject.Graph
	err = g.Provide(
		&inject.Object{Value: &calendar},
		&inject.Object{Value: &formBuilder},
		&inject.Object{Value: &guestbook},
		&inject.Object{Value: &handler},
		&inject.Object{Value: &ledger},
		&inject.Object{Value: &register},
		&inject.Object{Value: db},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err = g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Seed dates
	today, _ := date.Parse(time.Now())
	calendar.Add(today)
	calendar.Add(today.Add(1))

	// Start
	glog.Infoln("listening on", *flagHttp)
	glog.Fatal(http.ListenAndServe(*flagHttp, &handler))
}
