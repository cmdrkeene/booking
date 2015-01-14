package booking

import (
	"testing"
	"time"
)

func TestReservationManager(t *testing.T) {
	guest := guestId(1)
	feb1 := time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)
	feb2 := time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC)

	available := []time.Time{feb1, feb2}
	store := newReservationMemoryStore()
	manager := newReservationManager(available, store)

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
}
