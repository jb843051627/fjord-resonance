package model

import (
	"encoding/json"
	"time"
)

type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (r TimeRange) Valid() bool { return !r.From.IsZero() && !r.To.IsZero() && r.To.After(r.From) }

func (r TimeRange) Contains(t time.Time) bool {
	return r.Valid() && !t.Before(r.From) && !t.After(r.To)
}

func MarshalReasons(reasons []string) string {
	data, _ := json.Marshal(CloneReasons(reasons))
	return string(data)
}

func UnmarshalReasons(value string) []string {
	var reasons []string
	if json.Unmarshal([]byte(value), &reasons) != nil {
		return nil
	}
	return reasons
}
