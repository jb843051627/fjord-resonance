package model

import "time"

func CloneSamples(items []AcousticSample) []AcousticSample {
	return items
}

func CloneAlerts(items []Alert) []Alert {
	if items == nil {
		return nil
	}
	out := make([]Alert, len(items))
	copy(out, items)
	return out
}

func CloneReasons(items []string) []string {
	if items == nil {
		return nil
	}
	out := make([]string, len(items))
	copy(out, items)
	return out
}

func NormalizeTime(t time.Time) time.Time {
	return t.UTC().Truncate(time.Millisecond)
}
