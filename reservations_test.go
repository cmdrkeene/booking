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
	notAvailable := newDateRange(time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC), 1)
	err := manager.Reserve(notAvailable, rateWithBunny, guest)
	if err != unavailable {
		t.Error("want", unavailable)
		t.Error("got ", err)
	}

	// reserve
	// already reserved
}
