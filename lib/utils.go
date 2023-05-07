package lib

import (
	"time"
)

func GetDates(month time.Month, year int) (dates []string) {
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	for d := startOfMonth; !d.After(endOfMonth); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}
