package booking

import "testing"

func TestGuestbookFind(t *testing.T) {
	gb := newGuestbook()
	id, err := gb.Register("Brandon", "brandon@example.com")
	if err != nil {
		t.Error(err)
	}

	// not found
	_, err = gb.Find(guestId("999"))
	if err != notFound {
		t.Error("want", notFound)
		t.Error("got s", err)
	}

	// found
	record, err := gb.Find(id)
	if err != nil {
		t.Error(err)
	}

	if record.Email != "brandon@example.com" {
		t.Error("want", "brandon@example.com")
		t.Error("got ", record.Email)
	}
}

func TestGuestbookRegister(t *testing.T) {
	gb := newGuestbook()

	var tests = []struct {
		name  string
		email string
		err   error
	}{
		{"Brandon", "brandon@example.com", nil},
		{"Brandon X", "brandon@example.com", emailExists},
		{"", "brandon@example.com", nameMissing},
		{"Brandon", "", emailMissing},
		{"Brandon", "bork", emailInvalid},
	}

	for _, tt := range tests {
		guestId, err := gb.Register(tt.name, tt.email)
		if err != tt.err {
			t.Error("want", tt.err)
			t.Error("got ", err)
		}
		if err != nil {
			continue
		}
		if len(guestId) <= 0 {
			t.Error("want non zero length guestId")
			t.Error("got", guestId)
		}
	}
}
