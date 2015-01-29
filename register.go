package booking

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/golang/glog"
)

// Register is a container for booking records
type Register struct {
	Calendar *Calendar `inject:""`
	DB       *sql.DB   `inject:""`
}

const RegisterSchema = `
  CREATE TABLE Register (
    Checkin DATETIME UNIQUE NOT NULL REFERENCES Calendar,
    Checkout DATETIME UNIQUE NOT NULL REFERENCES Calendar,
    GuestId INTEGER NOT NULL REFERENCES Guestbook(Id),
    Id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    Rate TEXT NOT NULL,
    CONSTRAINT ck_Ckin_Less_Than_Ckout CHECK (Checkin <= Checkout)
  )
`

// Locator for a booking record
type bookingId uint8

func (id bookingId) String() string {
	return fmt.Sprintf("bookingId:%d", id)
}

type booking struct {
	Checkin  date.Date
	Checkout date.Date
	GuestId  guestId
	Id       bookingId
	Rate     rate
}

var (
	checkInAfterOut = errors.New("check in can't be after check out")
	stayTooShort    = errors.New("your stay is too short")
	unavailable     = errors.New("dates unavailable")
	minimumStay     = 1 // day
)

func (r *Register) Book(
	checkIn date.Date,
	checkOut date.Date,
	guest guestId,
	rate rate,
) (bookingId, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	bookingId, err := r.BookTx(tx, checkIn, checkOut, guest, rate)
	if err != nil {
		tx.Rollback()
		glog.Error(err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		glog.Error(err)
		return 0, err
	}

	return bookingId, nil
}

// Book with a passed Tx - caller is responsible for Commit/Rollback
func (r *Register) BookTx(
	tx *sql.Tx,
	checkIn date.Date,
	checkOut date.Date,
	guest guestId,
	rate rate,
) (bookingId, error) {
	// check in can't be after check out
	if checkIn.After(checkOut) {
		return 0, checkInAfterOut
	}

	// minimum stay length
	if checkIn.DaysApart(checkOut) < minimumStay {
		return 0, stayTooShort
	}

	// ensure on availablity calendar
	available, err := r.Calendar.Available(tx, checkIn, checkOut)
	if err != nil {
		glog.Error(err)
		return 0, err
	}
	if !available {
		return 0, unavailable
	}

	// TODO ensure can't be overbooked
	stmt, err := tx.Prepare(`
      insert into Register (Checkin, Checkout, GuestId, Rate)
      values ($1, $2, $3, $4)
    `)
	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(checkIn, checkOut, guest, rate)
	if err != nil {
		glog.Error(err)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
		return 0, err
	}

	return bookingId(lastId), nil
}

func (r *Register) Cancel(id bookingId) error {
	stmt, err := r.DB.Prepare(`delete from Register where Id=$1`)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}

func (r *Register) List() ([]booking, error) {
	stmt, err := r.DB.Prepare(`
    select Checkin, Checkout, GuestId, Id, Rate
    from Register
    order by CheckIn asc
  `)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		glog.Error(err)
		return []booking{}, err
	}
	defer rows.Close()

	var list []booking
	for rows.Next() {
		var checkin date.Date
		var checkout date.Date
		var guestId guestId
		var id bookingId
		var rate rate

		err := rows.Scan(
			&checkin,
			&checkout,
			&guestId,
			&id,
			&rate,
		)
		if err != nil {
			glog.Error(err)
			return []booking{}, err
		}

		list = append(list, booking{
			Checkin:  checkin,
			Checkout: checkout,
			GuestId:  guestId,
			Id:       id,
			Rate:     rate,
		})
	}

	return list, nil
}
