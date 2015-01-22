package booking

import "testing"

func TestAmount(t *testing.T) {
	a := amount(100)
	if s := a.String(); s != "$1.00" {
		t.Error("want $1.00")
		t.Error("got ", s)
	}
}
