package model

import "time"

func CloneSamples(items []AcousticSample) []AcousticSample {
	if items == nil {
		return nil
	}
	out := make([]AcousticSample, len(items))
	copy(out, items)
	return out
}

func CloneAlerts(items []Alert) []Alert {
	return items
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
