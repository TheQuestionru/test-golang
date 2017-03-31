package types

import (
	"time"

	"github.com/jinzhu/now"
)

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "2006-01-02T15:04:05Z" // RFC3339 without timezone
	dateFormatJson = `"` + DateFormat + `"`
)

// Date

type Date struct {
	time.Time
}

func ParseDate(s string) Date {
	t, _ := time.Parse(DateFormat, s)
	return Date{t}
}

func NewDate(t time.Time) Date {
	return Date{t.Truncate(time.Hour * 24)}
}

func Today() Date {
	return Date{now.BeginningOfDay()}
}

func (d Date) ToTime() Time {
	return Time{d.Time}
}

func (d Date) String() string {
	return d.Time.Format(DateFormat)
}

// Time

type Time struct {
	time.Time
}

func Now() Time {
	return Time{time.Now()}
}

func NewTime(t time.Time) Time {
	return Time{t}
}

func (t Time) ToNull() NullTime {
	return NewNullTime(t.Time)
}

func (t Time) String() string {
	return t.Time.Format(TimeFormat)
}

func (t Time) Truncate(d time.Duration) Time {
	t.Time = t.Time.Truncate(d)
	return t
}

func (t Time) Add(d time.Duration) Time {
	return Time{t.Time.Add(d)}
}

//
// NullTime
//

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func NullNow() NullTime {
	return NewNullTime(time.Now())
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{t, true}
}

func (t NullTime) String() string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	} else {
		return "null"
	}
}

func (t *NullTime) Set(time time.Time) {
	t.Time = time
	t.Valid = true
}

func (t *NullTime) Clear() {
	t.Time = time.Time{}
	t.Valid = false
}
