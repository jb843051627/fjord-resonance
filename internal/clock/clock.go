package clock

import "time"

type Clock interface{ Now() time.Time }

type System struct{}

func (System) Now() time.Time { return time.Now().UTC() }

type Fixed struct{ Value time.Time }

func (f Fixed) Now() time.Time { return f.Value }

func Today(c Clock) time.Time {
	now := c.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func SameInstant(a, b time.Time) bool { return a.UTC().Equal(b.UTC()) }
