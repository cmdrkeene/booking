package booking

import "testing"

func TestAmount(t *testing.T) {
	a := amount{25000}
	if a.InDollars() != "$250.00" {
		t.Error("want $250.00")
		t.Error("got ", a.InDollars())
	}
}
