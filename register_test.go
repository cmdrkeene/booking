package booking

import (
	"reflect"
	"testing"

	"github.com/facebookgo/inject"
)

func TestRegister(t *testing.T) {
	db := testDB()
	defer db.Close()
	var register Register
	var calendar Calendar
	err := inject.Populate(db, &register, &calendar)
	if err != nil {
		t.Error(err)
	}

	// book -> checkInAfterOut
	_, err = register.Book(
		newDate(2015, 1, 5),
		newDate(2015, 1, 2),
		guestId(123),
		withBunny,
	)
	if err != checkInAfterOut {
		t.Error("want", checkInAfterOut)
		t.Error("got ", err)
	}

	// book -> stayTooShort
	_, err = register.Book(
		newDate(2015, 1, 2),
		newDate(2015, 1, 2),
		guestId(123),
		withBunny,
	)
	if err != stayTooShort {
		t.Error("want", stayTooShort)
		t.Error("got ", err)
	}

	// book -> ok
	register.Calendar.Add(
		newDate(2015, 1, 2),
		newDate(2015, 1, 3),
		newDate(2015, 1, 4),
		newDate(2015, 1, 5),
	)
	id, err := register.Book(
		newDate(2015, 1, 2),
		newDate(2015, 1, 5),
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
			CheckIn:  newDate(2015, 1, 2),
			CheckOut: newDate(2015, 1, 5),
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
