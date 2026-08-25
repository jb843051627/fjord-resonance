package clock

import "time"

func RoundDown(value time.Time, step time.Duration) time.Time {
	if step <= 0 {
		return value
	}
	return value.Truncate(step)
}

func RoundUp(value time.Time, step time.Duration) time.Time {
	if step <= 0 {
		return value
	}
	down := value.Truncate(step)
	if down.Equal(value) {
		return value
	}
	return down.Add(step)
}

func IsBusinessHour(value time.Time, location *time.Location) bool {
	if location == nil {
		location = time.UTC
	}
	local := value.In(location)
	return local.Hour() >= 6 && local.Hour() < 18
}

func DaysBetween(from, to time.Time) int {
	return int(to.UTC().Truncate(24*time.Hour).Sub(from.UTC().Truncate(24*time.Hour)) / (24 * time.Hour))
}
