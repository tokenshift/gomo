package main

import (
	. "testing"
	"time"
)

func assertTimeEquals(t *T, expected, actual time.Time) bool {
	if expected.Equal(actual) {
		return true
	} else {
		t.Errorf("%v != %v", actual, expected)
		return false
	}
}

func TestRounding(t *T) {
	result := round(time.Date(2019, 5, 29, 11, 7, 14, 123, time.UTC), 15)
	assertTimeEquals(
		t,
		time.Date(2019, 5, 29, 11, 0, 0, 0, time.UTC),
		result)

	result = round(time.Date(2019, 5, 29, 11, 7, 29, 123, time.UTC), 15)
	assertTimeEquals(
		t,
		time.Date(2019, 5, 29, 11, 0, 0, 0, time.UTC),
		result)

	result = round(time.Date(2019, 5, 29, 11, 7, 30, 123, time.UTC), 15)
	assertTimeEquals(
		t,
		time.Date(2019, 5, 29, 11, 15, 0, 0, time.UTC),
		result)

	result = round(time.Date(2019, 5, 29, 11, 8, 14, 123, time.UTC), 15)
	assertTimeEquals(
		t,
		time.Date(2019, 5, 29, 11, 15, 0, 0, time.UTC),
		result)
}
