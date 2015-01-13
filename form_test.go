package booking

import (
	"net/http"
	"testing"
)

func TestBookingForm(t *testing.T) {
	var tests = []struct {
		url string
		err error
	}{
		{"/?name=", nameMissing},
		{"/?name=Brandon", emailMissing},
		{"/?name=Brandon&email=bork", emailInvalid},
		{"/?name=Brandon&email=a@b.com", phoneMissing},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212", datesMissing},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212&dates=bork", datesInvalid},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212&dates=2015-01-02&dates=2015-01-03", cardNumberMissing},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212&dates=2015-01-02&dates=2015-01-03&card_number=1111222233334444", cardMonthMissing},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212&dates=2015-01-02&dates=2015-01-03&card_number=1111222233334444&card_month=1", cardYearMissing},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212&dates=2015-01-02&dates=2015-01-03&card_number=1111222233334444&card_month=1&card_year=2016", cardCVCMissing},
		{"/?name=Brandon&email=a@b.com&phone=555-111-1212&dates=2015-01-02&dates=2015-01-03&card_number=1111222233334444&card_month=1&card_year=2016&card_cvc=111", nil},
	}

	for _, tt := range tests {
		req, err := http.NewRequest("POST", tt.url, nil)
		if err != nil {
			t.Error(err)
		}
		f, err := newBookingForm(req)
		if err != nil {
			if err != tt.err {
				t.Error("want", tt.err)
				t.Error("got ", err)
			}
			continue // don't check validation errors
		}

		if err := f.Validate(); tt.err != err {
			t.Error("want", tt.err)
			t.Error("got ", err)
		}
	}
}
