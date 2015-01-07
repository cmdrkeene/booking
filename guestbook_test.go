package booking

import "testing"

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
