package utils

// CronStringCreator returns a string
// for cron scheduler
func CronStringCreator(tz, m, h string) string {
	return "TZ=" + tz + " " + m + " " + h + " * * *"
}
