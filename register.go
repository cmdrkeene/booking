package booking

import "database/sql"

// Locator for a booking record
type bookingId uint32

// Official list of bookings
type Register struct {
	Calendar *Calendar `inject:""`
	DB       *sql.DB   `inject:""`
}

func (r *Register) Book(
	checkIn date,
	checkOut date,
	guest guestId,
	rate rate,
) (bookingId, error) {
	return bookingId(0), nil
}

func (r *Register) Cancel(bookingId) error {
	return nil
}
