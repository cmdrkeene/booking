package booking

import (
	"errors"
	"strings"
	"sync/atomic"
)

// Manages guest records
type Guestbook interface {
	Register(name, email string) (guestId, error)
	Find(guestId) (guest, error)
}

type guestId string

type guest struct {
	Email string
	Id    guestId
	Name  string
}

var (
	emailExists  = errors.New("email exists")
	emailMissing = errors.New("email missing")
	emailInvalid = errors.New("email invalid")
	nameMissing  = errors.New("name missing")
	notFound     = errors.New("not found")
)

type guestbook struct {
	guests     map[guestId]guest
	primaryKey uint32
}

func newGuestbook() guestbook {
	return guestbook{
		guests: make(map[guestId]guest),
	}
}

const minNameLength = 2

func (g guestbook) Register(name, email string) (guestId, error) {
	if name == "" {
		return "", nameMissing
	}

	if email == "" {
		return "", emailMissing
	}

	// use a crap check here and invalidate on failure to send
	if !strings.Contains(email, "@") {
		return "", emailInvalid
	}

	// check if email exists
	_, err := g.findByEmail(email)
	if err == nil {
		return "", emailExists
	}
	if err != notFound {
		return "", err
	}

	// add guest to database
	newGuest := guest{
		Name:  name,
		Email: email,
		Id:    g.newGuestId(),
	}
	g.guests[newGuest.Id] = newGuest

	return newGuest.Id, nil
}

func (g guestbook) Find(id guestId) (guest, error) {
	return guest{}, nil
}

func (g guestbook) findByEmail(s string) (guest, error) {
	for _, record := range g.guests {
		if record.Email == s {
			return record, nil
		}
	}
	return guest{}, notFound
}

func (g guestbook) newGuestId() guestId {
	id := atomic.AddUint32(&g.primaryKey, 1)
	return guestId(id)
}
