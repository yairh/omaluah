package calendar

import (
	"testing"
	"time"
)

func dayByHebrewDay(days []Day, hday int) Day {
	for _, d := range days {
		if d.HebrewDay == hday {
			return d
		}
	}
	return Day{}
}

func containsTitle(holidays []Holiday, title string) bool {
	for _, h := range holidays {
		if h.Title == title {
			return true
		}
	}
	return false
}

func TestHebDay(t *testing.T) {
	got := HebDay(time.Date(2026, time.September, 23, 0, 0, 0, 0, time.UTC))
	want := Day{
		Year:        2026,
		Month:       9,
		Day:         23,
		Date:        "2026-09-23",
		Weekday:     "Wednesday",
		HebrewYear:  5787,
		HebrewMonth: "Tishrei",
		HebrewDay:   12,
		HebrewDate:  "12 Tishrei 5787",
	}
	if !got.Equal(&want) {
		t.Errorf("HebDay error: got %+v, want: %+v", got, want)
	}
}

func TestDayEqual(t *testing.T) {
	a := HebDay(time.Date(2026, time.September, 23, 0, 0, 0, 0, time.UTC))
	b := HebDay(time.Date(2026, time.September, 24, 0, 0, 0, 0, time.UTC))
	if !a.Equal(&a) {
		t.Errorf("Equal should report the same Day as equal")
	}
	if a.Equal(&b) {
		t.Errorf("Equal should report different days as unequal")
	}
	if a.Equal(nil) || b.Equal(nil) {
		t.Errorf("Equal should report nil as unequal")
	}
}

func TestElulGrid(t *testing.T) {
	days, err := MonthGrid(2026, time.September, true)
	if err != nil {
		t.Fatalf("MonthGrid error: %v", err)
	}

	if len(days) != 29 {
		t.Errorf("expected 29 days in Elul 5786, got %d", len(days))
	}

	first := days[0]
	if first.HebrewYear != 5786 || first.HebrewMonth != "Elul" || first.HebrewDay != 1 {
		t.Errorf("unexpected first day: %+v", first)
	}
	if first.HebrewDate != "1 Elul 5786" {
		t.Errorf("unexpected HebrewDate: %q", first.HebrewDate)
	}
	if first.Date != "2026-08-14" {
		t.Errorf("unexpected Gregorian date: %q", first.Date)
	}
	if first.Weekday != "Friday" {
		t.Errorf("unexpected weekday: %q", first.Weekday)
	}

	last := days[len(days)-1]
	if last.HebrewDate != "29 Elul 5786" {
		t.Errorf("unexpected last HebrewDate: %q", last.HebrewDate)
	}
	if last.Date != "2026-09-11" {
		t.Errorf("unexpected last Gregorian date: %q", last.Date)
	}
}

func TestTishrei30Days(t *testing.T) {
	days, err := MonthGrid(2026, time.October, true)
	if err != nil {
		t.Fatalf("MonthGrid error: %v", err)
	}

	if len(days) != 30 {
		t.Errorf("expected 30 days in Tishrei 5787, got %d", len(days))
	}
	if days[0].HebrewDate != "1 Tishrei 5787" || days[0].Date != "2026-09-12" {
		t.Errorf("unexpected first day: %q / %q", days[0].HebrewDate, days[0].Date)
	}
	if days[29].HebrewDate != "30 Tishrei 5787" || days[29].Date != "2026-10-11" {
		t.Errorf("unexpected last day: %q / %q", days[29].HebrewDate, days[29].Date)
	}
}

func TestHolidayOverlay(t *testing.T) {
	days, err := MonthGrid(2026, time.September, true)
	if err != nil {
		t.Fatalf("MonthGrid error: %v", err)
	}

	first := dayByHebrewDay(days, 1)
	if !containsTitle(first.Holidays, "Rosh Chodesh Elul") {
		t.Errorf("expected Rosh Chodesh Elul on 1 Elul, got %+v", first.Holidays)
	}

	eve := dayByHebrewDay(days, 29)
	if !containsTitle(eve.Holidays, "Erev Rosh Hashana") {
		t.Errorf("expected Erev Rosh Hashana on 29 Elul, got %+v", eve.Holidays)
	}
}

func TestGridConsistency(t *testing.T) {
	days, err := MonthGrid(2026, time.November, true)
	if err != nil {
		t.Fatalf("MonthGrid error: %v", err)
	}

	for i, d := range days {
		if d.HebrewDay != i+1 {
			t.Errorf("day %d: HebrewDay out of order: %d", i, d.HebrewDay)
		}
		parsed, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			t.Fatalf("day %d: invalid Date %q: %v", i, d.Date, err)
		}
		if parsed.Year() != d.Year || int(parsed.Month()) != d.Month || parsed.Day() != d.Day {
			t.Errorf("day %d: Gregorian fields do not match Date %q", i, d.Date)
		}
		if d.Holidays != nil && len(d.Holidays) == 0 {
			t.Errorf("day %d: Holidays is empty but not nil", i)
		}
	}

	roshChodesh := dayByHebrewDay(days, 1)
	if !containsTitle(roshChodesh.Holidays, "Rosh Chodesh Cheshvan") {
		t.Errorf("expected Rosh Chodesh Cheshvan on 1 Cheshvan, got %+v", roshChodesh.Holidays)
	}
}

