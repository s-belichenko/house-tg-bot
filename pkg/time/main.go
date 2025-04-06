package time

import "time"

type ClockInterface interface {
	Now() time.Time
}

type Time struct{}

func (Time) Now() time.Time { return time.Now() }
