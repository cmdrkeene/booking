package booking

import "testing"

// full integration test
func TestEverything(t *testing.T) {
	// register
	// - send email to guest
	// reserve
	// - mark dateRange as reserved
	// - send email to manager
	// pay
	// - convert card to payment token
	// - save payment token in guestbook
	// - show cost before capture
	// - perform capture
	// - credit ledger with cost
	// - debit cost from ledger
	// confirm
	// - mark dateRange as booked
	// - send email to guest
	// - send email to manager
}
