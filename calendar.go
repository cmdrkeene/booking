package booking

import (
	"strings"
	"time"
)

type calendar map[time.Time]availability

type availability struct {
	Booked bool
}

func (c calendar) String() string {
	var lines []string
	lines = append(lines, "\n== Calendar ==")
	for t, a := range c {
		l := t.Format(dayFormat)
		if a.Booked {
			l = l + " (Booked)"
		}
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func (c calendar) SetAvailable(r dateRange) {
	for _, t := range r.Days() {
		c[t] = availability{Booked: false}
	}
}

func (c calendar) SetBooked(r dateRange) bool {
	// check if all available and not booked
	for _, t := range r.Days() {
		a, ok := c[t]
		if !ok {
			return false // unavailable
		}
		if a.Booked {
			return false // booked
		}
	}

	// mark it
	for _, t := range r.Days() {
		c[t] = availability{Booked: true}
	}
	return true
}
