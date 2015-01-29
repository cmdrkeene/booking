package booking

import "database/sql"

func testDB() *sql.DB {
	// connect
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// load schema
	var s Schema
	s.DB = db
	err = s.Load()
	if err != nil {
		panic(err)
	}

	return db
}
