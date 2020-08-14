package utils

import "time"

// GetCurrentDay returns number of current day
// Monday = 0, Tuesday = 1, Wednesday = 2, ...
func GetCurrentDay() int {
	currentDay := int(time.Now().Weekday())
	if currentDay == 0 {
		currentDay = 6
	} else {
		currentDay--
	}
	return currentDay
}
