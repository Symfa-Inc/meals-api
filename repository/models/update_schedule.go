package models

// UpdateSchedule request scheme
type UpdateSchedule struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	IsWorking *bool  `json:"isWorking"`
} //@name UpdateScheduleRequest
