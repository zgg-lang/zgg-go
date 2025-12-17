package utils

import (
	"sync/atomic"
	"time"
)

var timeDefaultLocal atomic.Pointer[time.Location]

func ParseTime(s, layout string, loc *time.Location) (t time.Time, unit string, err error) {
	if layout == "" {
		switch s {
		case "zero":
			return
		case "now":
			t = time.Now().In(loc)
			return
		case "today":
			t = time.Now().In(loc)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
			unit = "day"
			return
		case "yesterday":
			t = time.Now().In(loc).Add(-24 * time.Hour)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
			unit = "day"
			return
		case "tomorrow":
			t = time.Now().In(loc).Add(24 * time.Hour)
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
			unit = "day"
			return
		}
		switch len(s) {
		case 6:
			layout = "060102"
			unit = "day"
		case 8:
			layout = "20060102"
			unit = "day"
		case 10:
			layout = "2006-01-02"
			unit = "day"
		case 14:
			layout = "20060102150405"
			unit = "second"
		case 19:
			layout = "2006-01-02 15:04:05"
			unit = "second"
		case 19 + 3:
			layout = "2006-01-02 15:04:05-07"
			unit = "second"
		case 19 + 5:
			layout = "2006-01-02 15:04:05-0700"
			unit = "second"
		case 19 + 4:
			layout = "2006-01-02 15:04:05.000"
		case 19 + 4 + 3:
			layout = "2006-01-02 15:04:05.000-07"
		case 19 + 4 + 5:
			layout = "2006-01-02 15:04:05.000-0700"
		}
	}
	if loc == nil {
		loc = timeDefaultLocal.Load()
	}
	t, err = time.ParseInLocation(layout, s, loc)
	return
}

func GetDefaultTimeLocal() *time.Location {
	return timeDefaultLocal.Load()
}

func SetDefaultTimeLocal(loc *time.Location) {
	timeDefaultLocal.Store(loc)
}

func init() {
	SetDefaultTimeLocal(time.Local)
}
