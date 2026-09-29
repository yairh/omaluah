package calendar

import (
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/hebcal"
)

const layout = "2006-01-02"

type Holiday struct {
	Title       string   `json:"title"`
	HebrewTitle string   `json:"hebrewTitle"`
	Categories  []string `json:"categories"`
}

type Day struct {
	Year        int       `json:"year"`
	Month       int       `json:"month"`
	Day         int       `json:"day"`
	Date        string    `json:"date"`
	Weekday     string    `json:"weekday"`
	HebrewYear  int       `json:"hebrewYear"`
	HebrewMonth string    `json:"hebrewMonth"`
	HebrewDay   int       `json:"hebrewDay"`
	HebrewDate  string    `json:"hebrewDate"`
	Holidays    []Holiday `json:"holidays"`
}

func (a *Day) Equal(b *Day) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Year != b.Year || a.Month != b.Month || a.Day != b.Day {
		return false
	}
	if a.Date != b.Date || a.Weekday != b.Weekday {
		return false
	}
	if a.HebrewYear != b.HebrewYear || a.HebrewMonth != b.HebrewMonth ||
		a.HebrewDay != b.HebrewDay || a.HebrewDate != b.HebrewDate {
		return false
	}
	if len(a.Holidays) != len(b.Holidays) {
		return false
	}
	for i := range a.Holidays {
		ha, hb := a.Holidays[i], b.Holidays[i]
		if ha.Title != hb.Title || ha.HebrewTitle != hb.HebrewTitle {
			return false
		}
		if len(ha.Categories) != len(hb.Categories) {
			return false
		}
		for j := range ha.Categories {
			if ha.Categories[j] != hb.Categories[j] {
				return false
			}
		}
	}
	return true
}

// MonthGrid returns a flat array with every day (1..29/30) of the Hebrew month
// that contains the first day of the given Gregorian year/month. Holidays from
// the Hebcal calendar are overlaid on each day they occur.
func MonthGrid(year int, month time.Month, il bool) ([]Day, error) {
	firstOfGregMonth := hdate.FromGregorian(year, month, 1)
	hyear := firstOfGregMonth.Year()
	hmonth := firstOfGregMonth.Month()

	daysInMonth := hdate.DaysInMonth(hmonth, hyear)
	gridStart := hdate.New(hyear, hmonth, 1)
	gridEnd := hdate.New(hyear, hmonth, daysInMonth)

	options := hebcal.CalOptions{Start: gridStart, End: gridEnd, IL: il}
	res, err := hebcal.HebrewCalendar(&options)
	if err != nil {
		return nil, err
	}

	holidayMap := map[string][]Holiday{}
	for _, el := range res {
		d := el.GetDate().Gregorian().Format(layout)
		holidayMap[d] = append(holidayMap[d], Holiday{
			Title:       el.Render("en"),
			HebrewTitle: el.Render("he"),
			Categories:  el.GetCategories(),
		})
	}

	grid := make([]Day, 0, daysInMonth)
	for d := 1; d <= daysInMonth; d++ {
		hd := hdate.New(hyear, hmonth, d)
		greg := hd.Gregorian()
		dateStr := greg.Format(layout)
		grid = append(grid, Day{
			Year:        greg.Year(),
			Month:       int(greg.Month()),
			Day:         greg.Day(),
			Date:        dateStr,
			Weekday:     hd.Weekday().String(),
			HebrewYear:  hyear,
			HebrewMonth: hd.MonthName("en"),
			HebrewDay:   d,
			HebrewDate:  hd.String(),
			Holidays:    holidayMap[dateStr],
		})
	}

	return grid, nil
}

func HebDay(greg time.Time) Day {
	hd := hdate.FromGregorian(greg.Year(), greg.Month(), greg.Day())
	dateStr := greg.Format(layout)
	holidays := []Holiday{}
	for _, el := range hebcal.GetHolidaysOnDate(hd, true) {
		holidays = append(holidays, Holiday{
			Title:       el.Render("en"),
			HebrewTitle: el.Render("he"),
			Categories:  el.GetCategories(),
		})
	}

	return Day{
		Year:        greg.Year(),
		Month:       int(greg.Month()),
		Day:         greg.Day(),
		Date:        dateStr,
		Weekday:     hd.Weekday().String(),
		HebrewYear:  hd.Year(),
		HebrewMonth: hd.MonthName("en"),
		HebrewDay:   hd.Day(),
		HebrewDate:  hd.String(),
		Holidays:    holidays,
	}
}
