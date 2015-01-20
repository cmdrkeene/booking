package booking

import (
	"net/http"
	"testing"
	"time"
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

func TestBookingFormQueryString(t *testing.T) {
	form := bookingForm{}
	form.Name = "Brandon"
	form.Email = "a@b.com"
	form.Phone = "(555) 111-1212"
	form.Dates = []time.Time{
		time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC),
	}
	form.CardNumber = "1111222233334444"
	form.CardMonth = "1"
	form.CardYear = "2015"
	form.CardCVC = "111"
	want := "card_cvc=111&card_month=1&card_number=1111222233334444&card_year=2015&dates=2015-02-01&dates=2015-02-02&email=a%40b.com&name=Brandon&phone=%28555%29+111-1212"
	if s := form.Encode(); s != want {
		t.Error("want", want)
		t.Error("got ", s)
	}
}
