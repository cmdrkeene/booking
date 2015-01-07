package booking

// Manages guest records
type Guestbook interface {
	Register(name, email string) (guestId, error)
	Find(guestId) (guest, error)
}

type guestId string

// The rando sleeping/having-sex in your bed
type guest struct {
	Email string
	Id    guestId
	Name  string
}

type guestbook struct {
	guests []guest
}

func (g guestbook) Register(name, email string) (guestId, error) {
	return "", nil
}

func (g guestbook) Find(id guestId) (guest, error) {
	return guest{}, nil
}
