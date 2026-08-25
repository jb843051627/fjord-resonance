package clock

import "time"

type Window struct{ Start, End time.Time }

func DayWindow(c Clock, day time.Time) Window {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return Window{Start: start, End: start.Add(24 * time.Hour)}
}

func Overlaps(a, b Window) bool { return a.Start.Before(b.End) && b.Start.Before(a.End) }

func Contains(window Window, instant time.Time) bool {
	return !instant.Before(window.Start) && instant.Before(window.End)
}

func Expand(window Window, padding time.Duration) Window {
	return Window{Start: window.Start.Add(-padding), End: window.End.Add(padding)}
}
