package main

import (
	"math"
	"time"
)

func round(t time.Time, minutes int) time.Time {
	if minutes >= 0 {
		rounded := float64(t.Minute()) + float64(t.Second())/60
		rounded = rounded / float64(minutes)
		rounded = math.Round(rounded)
		rounded = rounded * float64(minutes)

		return time.Date(
			t.Year(), t.Month(), t.Day(),
			t.Hour(), int(rounded), 0, 0, t.Location())
	} else {
		return t
	}
}
