package booking

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

type guestId uint32

func (id guestId) String() string {
	return fmt.Sprintf("guest:%d", id)
}

type guest struct {
	Email string
	Id    guestId
	Name  string
}

type registrar interface {
	Register(name, email string) (guestId, error)
}

type guestFinder interface {
	Find(guestId) (guest, error)
}

var (
	emailExists  = errors.New("email exists")
	emailInvalid = errors.New("email invalid")
	emailMissing = errors.New("email missing")
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
		return 0, nameMissing
	}

	if email == "" {
		return 0, emailMissing
	}

	// use a crap check here and invalidate on failure to send
	if !strings.Contains(email, "@") {
		return 0, emailInvalid
	}

	// check if email exists
	_, err := g.findByEmail(email)
	if err == nil {
		return 0, emailExists
	}
	if err != notFound {
		return 0, err
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
	if record, ok := g.guests[id]; ok {
		return record, nil
	}

	return guest{}, notFound
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
