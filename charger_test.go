package booking

import "testing"

func TestFakeCharger(t *testing.T) {
	var c charger
	c = fakeCharger{}
	err := c.Charge(creditCard{}, amount{})
	if err != nil {
		t.Error(err)
	}
}
