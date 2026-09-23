package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"omaluah/calendar"
)

func main() {
	now := time.Now()
	today := flag.Bool("today", false, "Return today's hebrew date")
	year := flag.Int("year", now.Year(), "Gregorian Year")
	month := flag.Int("month", int(now.Month()), "Gregorian Month")
	il := flag.Bool("il", true, "Israeli Schedule")
	flag.Parse()

	if *today == true {
		day := calendar.HebDay(now)
		if err := json.NewEncoder(os.Stdout).Encode(day); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	days, err := calendar.MonthGrid(*year, time.Month(*month), *il)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if err := json.NewEncoder(os.Stdout).Encode(days); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

