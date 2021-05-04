package handlers

import (
	"time"
)

type Timeframe string

const (
	Day   Timeframe = "day"
	Month Timeframe = "month"
)

func (t Timeframe) String() string {
	return string(t)
}

const (
	DayDuration   time.Duration = time.Hour * 24
	MonthDuration time.Duration = DayDuration * 30
)
