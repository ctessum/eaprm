package eaprm

import (
	"reflect"
	"testing"
	"time"
)

func TestHour(t *testing.T) {
	times := []time.Time{
		time.Date(1969, time.July, 20, 23, 0, 0, 0, time.UTC),
		time.Date(1929, time.October, 24, 3, 0, 0, 0, time.UTC),
		time.Date(1066, time.October, 14, 11, 0, 0, 0, time.UTC),
	}
	hours := Hour(times)
	sum := 0.0
	for _, h := range hours {
		for _, hh := range h {
			sum += hh
		}
	}
	if sum != 3 {
		t.Errorf("sum: have %g but want %d", sum, 3)
	}
	if hours[23][0] != 1 {
		t.Error("times[23][0] should equal 1")
	}
	if hours[3][1] != 1 {
		t.Error("times[3][1] should equal 1")
	}
	if hours[11][2] != 1 {
		t.Error("times[11][2] should equal 1")
	}
}

func TestMonth(t *testing.T) {
	times := []time.Time{
		time.Date(1969, time.July, 20, 23, 0, 0, 0, time.UTC),
		time.Date(1929, time.October, 24, 3, 0, 0, 0, time.UTC),
		time.Date(1066, time.October, 14, 11, 0, 0, 0, time.UTC),
	}
	months := Month(times)
	sum := 0.0
	for _, m := range months {
		for _, mm := range m {
			sum += mm
		}
	}
	if sum != 3 {
		t.Errorf("sum: have %g but want %d", sum, 3)
	}
	if months[6][0] != 1 {
		t.Error("times[6][0] should equal 1")
	}
	if months[9][1] != 1 {
		t.Error("times[9][1] should equal 1")
	}
	if months[9][2] != 1 {
		t.Error("times[9][2] should equal 1")
	}
}

func TestWeekend(t *testing.T) {
	times := []time.Time{
		time.Date(1969, time.July, 20, 23, 0, 0, 0, time.UTC),
		time.Date(1929, time.October, 24, 3, 0, 0, 0, time.UTC),
		time.Date(1066, time.October, 14, 11, 0, 0, 0, time.UTC),
	}
	weekend := Weekend(times)
	want := []float64{1, 0, 1}
	if !reflect.DeepEqual(weekend, want) {
		t.Errorf("weekend: have %v, want %v", weekend, want)
	}
}

func TestYear(t *testing.T) {
	times := []time.Time{
		time.Date(2000, time.July, 20, 23, 0, 0, 0, time.UTC),
		time.Date(2001, time.October, 24, 3, 0, 0, 0, time.UTC),
		time.Date(2003, time.October, 14, 11, 0, 0, 0, time.UTC),
	}
	years, err := Year(times, 2000, 2003)
	if err != nil {
		t.Fatal(err)
	}
	sum := 0.0
	for _, m := range years {
		for _, mm := range m {
			sum += mm
		}
	}
	if sum != 3 {
		t.Errorf("sum: have %g but want %d", sum, 3)
	}
	if years[0][0] != 1 {
		t.Error("times[0][0] should equal 1")
	}
	if years[1][1] != 1 {
		t.Error("times[1][1] should equal 1")
	}
	if years[3][2] != 1 {
		t.Error("times[3][2] should equal 1")
	}
}
