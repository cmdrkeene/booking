package booking

import (
	"testing"
	"time"
)

func TestReservationManager(t *testing.T) {
	guest := guestId(1)
	feb1 := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)

	availability := &testAvailablity{[]time.Time{feb1, feb1.Add(day)}}
	reservations := newReservationMemoryStore()
	manager := newReservationManager(availability, reservations)

	// not available
	err := manager.Reserve(newDateRange(feb1, 7), rateWithBunny, guest)
	if err != unavailable {
		t.Error("want", unavailable)
		t.Error("got ", err)
	}

	// reserve
	err = manager.Reserve(newDateRange(feb1, 2), rateWithBunny, guest)
	if err != nil {
		t.Error("want nil")
		t.Error("got  ", err)
	}

	// already reserved
	err = manager.Reserve(newDateRange(feb1, 2), rateWithBunny, guest)
	if err != unavailable {
		t.Error("want", unavailable)
		t.Error("got ", err)
	}

	// reserve newly available
	feb3 := feb1.Add(2 * day)
	availability.Add(feb3)
	err = manager.Reserve(newDateRange(feb3, 1), rateWithBunny, guest)
	if err != nil {
		t.Error("want nil")
		t.Error("got  ", err)
	}
}

type testAvailablity struct {
	list []time.Time
}

func (a *testAvailablity) Add(t time.Time) error {
	a.list = append(a.list, t)
	return nil
}

func (a *testAvailablity) Remove(t time.Time) error {
	for i, e := range a.list {
		if t.Equal(e) {
			a.list = append(a.list[:i], a.list[i+1:]...)
			break
		}
	}
	return nil
}

func (a *testAvailablity) List() ([]time.Time, error) {
	return a.list, nil
}
