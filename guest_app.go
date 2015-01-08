package booking

/*
GuestApp is a web app for guests to register, make bookings, etc.

Pages:
* Home - shows a calendar of availability
* Book - date picker, if unavailable show error and repeat
* Register - after date is picked, register as a guest
* Pay - after dates picked, guest registered, charge credit card
* Confirmed - after credit card charged, confirm reservation
*/
type guestApp struct {
	charger   charger   // charges cards
	registrar registrar // registers guests
	reserver  reserver  // reserves dates
}

func newGuestApp(chg charger, reg registrar, res reserver) guestApp {
	return guestApp{}
}
