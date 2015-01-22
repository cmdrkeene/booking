package booking

import (
	"reflect"
	"testing"
)

func TestGuestbook(t *testing.T) {
	var guestbook Guestbook
	guestbook.DB = testDB()
	defer guestbook.DB.Close()

	var err error

	// Lookup -> guestNotFound
	_, err = guestbook.Lookup(guestId(1))
	if err != guestNotFound {
		t.Error("want", guestNotFound)
		t.Error("got ", err)
	}

	// Register -> ok
	name, _ := newName("B K")
	email, _ := newEmail("a@b.com")
	phoneNumber, _ := newPhoneNumber("555-111-2222")
	id, err := guestbook.Register(name, email, phoneNumber)
	if err != nil {
		t.Error(err)
	}

	// Lookup -> ok
	found, err := guestbook.Lookup(id)
	if err != nil {
		t.Error(err)
	}

	guest := guest{
		Email:       email,
		Id:          id,
		Name:        name,
		PhoneNumber: phoneNumber,
	}

	if !reflect.DeepEqual(guest, found) {
		t.Error("want", guest)
		t.Error("got ", found)
	}

	guestbook.DB.Close()

	// Lookup -> err db closed
	_, err = guestbook.Lookup(id)
	if err == nil {
		t.Error("expected error")
	}
}

func TestGuestId(t *testing.T) {
	id := guestId(123)
	if s := id.String(); s != "guestId:123" {
		t.Error("want guestId:123")
		t.Error("got ", s)
	}
}

func TestEmail(t *testing.T) {
	var tests = []struct {
		input  string
		output string
		err    error
	}{
		{"", "", invalidEmail},
		{" ", "", invalidEmail},
		{"a", "", invalidEmail},
		{"@", "", invalidEmail},
		{"a@", "", invalidEmail},
		{"@b", "", invalidEmail},
		{"a@b", "a@b", nil},
		{" a@b", "a@b", nil},
		{"a@b ", "a@b", nil},
		{" a@b ", "a@b", nil},
		{"user@example.com", "user@example.com", nil},
	}

	for _, tt := range tests {
		email, err := newEmail(tt.input)
		if err != tt.err {
			t.Error("want", tt.err)
			t.Error("got ", err)
			continue
		}

		if !email.Equal(tt.output) {
			t.Error("want", tt.output)
			t.Error("got ", email)
		}
	}
}

func TestGuestName(t *testing.T) {
	var tests = []struct {
		input  string
		output string
		err    error
	}{
		{"", "", invalidName},
		{" ", "", invalidName},
		{"B", "", invalidName},
		{"B K", "B K", nil},
		{" B K ", "B K", nil},
		{"Brandon Allen Keene", "Brandon Allen Keene", nil},
	}

	for _, tt := range tests {
		name, err := newName(tt.input)
		if tt.err != err {
			t.Error("want", tt.err)
			t.Error("got ", err)
			continue
		}

		if !name.Equal(tt.output) {
			t.Error("want", tt.output)
			t.Error("got ", name)
		}
	}
}
