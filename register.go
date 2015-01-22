package booking

import (
	"database/sql"
	"errors"

	"github.com/golang/glog"
)

// Locator for a booking record
type bookingId uint32

type booking struct {
	CheckIn  date
	CheckOut date
	GuestId  guestId
	Id       bookingId
	Rate     rate
}

// Official list of bookings
type Register struct {
	Calendar *Calendar `inject:""`
	DB       *sql.DB   `inject:""`

	book         *sql.Stmt
	cancel       *sql.Stmt
	list         *sql.Stmt
	tableCreated bool
}

func (r *Register) createTableOnce() {
	if r.tableCreated {
		return
	}

	_, err := r.DB.Exec(
		`create table Register (
      CheckIn datetime not null, 
      CheckOut datetime not null, 
      GuestId integer not null, 
      Id integer primary key autoincrement not null, 
      Rate text not null
    )`,
	)
	if err == nil {
		r.tableCreated = true
		glog.Info("Register table created")
	} else {
		glog.Warning(err)
	}
}

var (
	checkInAfterOut = errors.New("check in can't be after check out")
	stayTooShort    = errors.New("your stay is too short")
	unavailable     = errors.New("unavailable")
	minimumStay     = 1 // day
)

func (r *Register) Book(
	checkIn date,
	checkOut date,
	guest guestId,
	rate rate,
) (bookingId, error) {
	r.createTableOnce()

	// check in can't be after check out
	if checkIn.After(checkOut) {
		return 0, checkInAfterOut
	}

	// minimum stay length
	if checkIn.DaysApart(checkOut) < minimumStay {
		return 0, stayTooShort
	}

	// ensure on availablity calendar
	available, err := r.Calendar.Available(checkIn, checkOut)
	if err != nil {
		glog.Error(err)
		return 0, err
	}
	if !available {
		return 0, unavailable
	}

	// TODO ensure can't be overbooked

	if r.book == nil {
		var err error
		r.book, err = r.DB.Prepare(`
      insert into Register (CheckIn, CheckOut, GuestId, Rate)
      values ($1, $2, $3, $4)
    `)
		if err != nil {
			panic(err)
		}
	}

	result, err := r.book.Exec(checkIn, checkOut, guest, rate)
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
	r.createTableOnce()

	if r.cancel == nil {
		var err error
		r.cancel, err = r.DB.Prepare(`delete from Register where Id=$1`)
		if err != nil {
			panic(err)
		}
	}

	_, err := r.cancel.Exec(id)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}

func (r *Register) List() ([]booking, error) {
	r.createTableOnce()

	if r.list == nil {
		var err error
		r.list, err = r.DB.Prepare(`
      select CheckIn, CheckOut, GuestId, Id, Rate
      from Register
      order by CheckIn asc
    `)
		if err != nil {
			panic(err)
		}
	}

	rows, err := r.list.Query()
	if err != nil {
		glog.Error(err)
		return []booking{}, err
	}
	defer rows.Close()

	var list []booking
	for rows.Next() {
		var checkIn date
		var checkOut date
		var guestId guestId
		var id bookingId
		var rate rate

		err := rows.Scan(
			&checkIn,
			&checkOut,
			&guestId,
			&id,
			&rate,
		)
		if err != nil {
			glog.Error(err)
			return []booking{}, err
		}

		list = append(list, booking{
			CheckIn:  checkIn,
			CheckOut: checkOut,
			GuestId:  guestId,
			Id:       id,
			Rate:     rate,
		})
	}

	return list, nil
}
