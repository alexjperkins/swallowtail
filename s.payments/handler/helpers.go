package handler

import "time"

// CurrentMonthStartTimestamp returns the timestamp of the start of the current month.
// This is defined as the 1st of every month at 00:00:00
func currentMonthStartTimestamp() time.Time {
	now := time.Now().Truncate(time.Hour)
	daysIntoMonth := now.Day()
	return now.AddDate(0, 0, -daysIntoMonth)
}
