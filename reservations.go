package booking

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var unavailable = errors.New("unavailable")

type reservationId uint32

func (id reservationId) String() string {
	return fmt.Sprintf("reservation:%d", id)
}

type reservation struct {
	id    reservationId
	guest guestId
	dates dateRange
	rate  rateCode
}

type reservationStore interface {
	Reserve(dateRange, guestId, rateCode) (reservationId, error)
	Cancel(reservationId) error
	List() ([]reservation, error)
}

type reservationTable struct {
	reserve *sql.Stmt
	cancel  *sql.Stmt
	list    *sql.Stmt
}

func newReservationTable(db *sql.DB) reservationTable {
	t := reservationTable{}

	_, err := db.Exec(`
    CREATE TABLE Reservation (
      Cancelled bool DEFAULT false NOT NULL,
      End datetime NOT NULL,
      GuestId int NOT NULL,
      Id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
      RateCode INTEGER NOT NULL,
      Start datetime NOT NULL
    )
  `)
	if err != nil {
		panic(err)
	}

	// var err error

	t.reserve, err = db.Prepare(`
    INSERT INTO Reservation
    (Start, End, GuestId, RateCode)
    VALUES
    ($1, $2, $3, $4)
  `)
	if err != nil {
		panic(err)
	}

	t.list, err = db.Prepare(`
    SELECT
    Start, End, GuestId, RateCode
    FROM Reservation
  `)
	if err != nil {
		panic(err)
	}

	return t
}

func (table reservationTable) Reserve(
	dates dateRange,
	guest guestId,
	rate rateCode,
) (reservationId, error) {
	res, err := table.reserve.Exec(
		dates.Start(),
		dates.End(),
		guest,
		rate,
	)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return reservationId(lastId), nil
}

func (table reservationTable) Cancel(reservationId) error {
	return nil
}

func (table reservationTable) List() ([]reservation, error) {
	rows, err := table.list.Query()
	if err != nil {
		return []reservation{}, err
	}
	defer rows.Close()

	var list []reservation
	for rows.Next() {
		var start time.Time
		var end time.Time
		var guest int
		var rate rateCode
		err = rows.Scan(&start, &end, &guest, &rate)
		if err != nil {
			return []reservation{}, err
		}

		rec := reservation{
			dates: newDateRangeBetween(start, end),
			guest: guestId(guest),
			rate:  rateCode(rate),
		}
		list = append(list, rec)
	}

	err = rows.Err()
	if err != nil {
		return []reservation{}, err
	}

	return list, nil
}

type reserver interface {
	Available() []time.Time
	IsAvailable(dateRange) bool
	Reserve(dateRange, rateCode, guestId) error
}

type reservationManager struct {
	availability availabilityStore
	reservations reservationStore
}

func newReservationManager(a availabilityStore, r reservationStore) reservationManager {
	return reservationManager{a, r}
}

func (m reservationManager) Reserve(dates dateRange, rate rateCode, guest guestId) error {
	// check available
	available, err := m.availability.List()
	if err != nil {
		return err
	}
	if !dates.Coincident(available) {
		return unavailable
	}

	// check reserved
	reserved, err := m.reserved()
	if err != nil {
		return err
	}
	if dates.Coincident(reserved) {
		return unavailable
	}

	// save
	reservationId, err := m.reservations.Reserve(dates, guest, rate)
	if err != nil {
		return err
	}

	log.Println(reservationId, "created by", guest, "for", dates)
	return nil
}

func (m reservationManager) reserved() ([]time.Time, error) {
	list, err := m.reservations.List()
	if err != nil {
		return []time.Time{}, err
	}

	var times []time.Time
	for _, r := range list {
		for _, t := range r.dates.EachDay() {
			times = append(times, t)
		}
	}
	return times, nil
}
