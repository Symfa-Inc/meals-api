package utils

import "time"

func GetCurrentDay() int {
	currentDay := int(time.Now().Weekday())
	if currentDay == 0 {
		currentDay = 6
	} else {
		currentDay -= 1
	}
	return currentDay
}
