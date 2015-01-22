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
	var name name
	name.Set("B K")

	var email email
	email.Set("a@b.com")

	phoneNumber := phoneNumber("555-111-2222")

	id, err := guestbook.Register(&name, &email, phoneNumber)
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
		ok     bool
	}{
		{"", "", false},
		{" ", "", false},
		{"a", "", false},
		{"@", "", false},
		{"a@", "", false},
		{"@b", "", false},
		{"a@b", "a@b", true},
		{" a@b", "a@b", true},
		{"a@b ", "a@b", true},
		{" a@b ", "a@b", true},
		{"user@example.com", "user@example.com", true},
	}

	for _, tt := range tests {
		var email email

		if ok := email.Set(tt.input); ok != tt.ok {
			t.Error("want", tt.ok)
			t.Error("got ", ok)
			continue
		}

		if !email.Equal(tt.output) {
			t.Error("want", tt.output)
			t.Error("got ", email)
		}
	}
}

func TestUserName(t *testing.T) {
	var tests = []struct {
		input  string
		output string
		ok     bool
	}{
		{"", "", false},
		{" ", "", false},
		{"B", "", false},
		{"B K", "B K", true},
		{" B K ", "B K", true},
	}

	for _, tt := range tests {
		var name name

		if ok := name.Set(tt.input); ok != tt.ok {
			t.Error("want", tt.ok)
			t.Error("got ", ok)
			continue
		}

		if !name.Equal(tt.output) {
			t.Error("want", tt.output)
			t.Error("got ", name)
		}
	}
}
