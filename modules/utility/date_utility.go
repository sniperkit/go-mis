package utility

import (
	"log"
	"time"
)

const (
	defaultFormat = "2006-01-02"
)

// GetCurrentDateToString - Get only date to string
func GetCurrentDateToString() string {
	return time.Now().Local().Format(defaultFormat)
}

// StringToDate - Get date from string
func StringToDate(strDate string) (time.Time, error) {
	return time.Parse(defaultFormat, strDate)
}

// IsBeforeToday - Compare whether param date is before current date or not
func IsBeforeToday(dateParam time.Time) bool {
	now, err := StringToDate(GetCurrentDateToString())
	if err != nil {
		log.Println("#ERROR IsAfterNow: Failed to compare date")
		return false
	}
	return dateParam.Before(now)
}

// IsAfterToday - Compare whether param date is after current date or not
func IsAfterToday(dateParam time.Time) bool {
	now, err := StringToDate(GetCurrentDateToString())
	if err != nil {
		log.Println("#ERROR IsAfterNow: Failed to compare date")
		return false
	}
	return dateParam.After(now)
}

// IsToday - Check if dateParam is today or not
func IsToday(dateParam time.Time) bool {
	now, err := StringToDate(GetCurrentDateToString())
	if err != nil {
		log.Println("#ERROR IsAfterNow: Failed to compare date")
		return false
	}
	return !now.After(dateParam) && !now.Before(dateParam)
}

// IsYesterday - Check whether the date param is yesterday
func IsYesterday(date time.Time) bool {
	yesterday := time.Now().Add(-24 * time.Hour)
	return date == yesterday
}
