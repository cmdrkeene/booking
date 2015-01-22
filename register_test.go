package booking

import (
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	var register Register
	register.DB = testDB()
	var calendar Calendar
	calendar.DB = register.DB
	register.Calendar = &calendar
	defer register.DB.Close()

	// book
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
