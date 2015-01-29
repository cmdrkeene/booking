package booking

import (
	"reflect"
	"testing"

	"github.com/cmdrkeene/booking/pkg/date"
	"github.com/facebookgo/inject"
)

func TestRegister(t *testing.T) {
	db := testDB()
	defer db.Close()
	var calendar Calendar
	var register Register
	err := inject.Populate(&calendar, db, &register)
	if err != nil {
		t.Error(err)
	}

	// book -> checkInAfterOut
	_, err = register.Book(
		date.New(2015, 1, 5),
		date.New(2015, 1, 2),
		guestId(123),
		withBunny,
	)
	if err != checkInAfterOut {
		t.Error("want", checkInAfterOut)
		t.Error("got ", err)
	}

	// book -> stayTooShort
	_, err = register.Book(
		date.New(2015, 1, 2),
		date.New(2015, 1, 2),
		guestId(123),
		withBunny,
	)
	if err != stayTooShort {
		t.Error("want", stayTooShort)
		t.Error("got ", err)
	}

	// book -> unavailable
	_, err = register.Book(
		date.New(2015, 3, 1),
		date.New(2015, 3, 20),
		guestId(123),
		withBunny,
	)
	if err != unavailable {
		t.Error("want", unavailable)
		t.Error("got ", err)
	}

	// book -> ok
	calendar.Add(
		date.New(2015, 1, 2),
		date.New(2015, 1, 3),
		date.New(2015, 1, 4),
		date.New(2015, 1, 5),
	)

	id, err := register.Book(
		date.New(2015, 1, 2),
		date.New(2015, 1, 5),
		guestId(123),
		withBunny,
	)
	if err != nil {
		t.Error(err)
	}

	// list
	list, err := register.List()
	if err != nil {
		t.Error(err)
	}
	want := []booking{
		booking{
			Checkin:  date.New(2015, 1, 2),
			Checkout: date.New(2015, 1, 5),
			GuestId:  guestId(123),
			Id:       id,
			Rate:     withBunny,
		},
	}
	if !reflect.DeepEqual(want, list) {
		t.Error("want", want)
		t.Error("got ", list)
	}

	// cancel
	err = register.Cancel(id)
	if err != nil {
		t.Error(err)
	}

	list, err = register.List()
	if err != nil {
		t.Error(err)
	}
	if l := len(list); l != 0 {
		t.Error("want 0")
		t.Error("got ", l)
	}
}
