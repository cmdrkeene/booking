package booking

type guestId string

// The rando sleeping/having-sex in your bed
type guest struct {
	Name          string
	Email         string
	Id            guestId
	PaymentTokens []paymentToken
}

type guestbook struct {
	guests []guest
}

func (g guestbook) Register(name, email string) (guest, error) {
	return guest{}, nil
}

// trades creditCard for paymentToken and
func (g guestbook) AddCard(guest, creditCard) error {
	return nil
}
