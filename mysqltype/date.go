package mysqltype

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"time"
)

// Date support MySQL Date type
// https://dev.mysql.com/doc/refman/8.0/en/datetime.html
type Date struct {
	src time.Time `gorm:"type:date"`
}

// NewDate Create new Date from time.Date
func NewDate(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return NewDateFromTime(t)
}

// NewDateFromTime Create new Date from Time
func NewDateFromTime(t time.Time) Date {
	return Date{src: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

// MinDate Minimum Date
func MinDate() Date {
	return NewDateFromTime(time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// MaxDate Maximum Date
func MaxDate() Date {
	return NewDateFromTime(time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC))
}

// NowDate Create Now time for MySQL DataBase
func NowDate() Date {
	return NewDateFromTime(time.Now())
}

// After behavior as time.Time
func (dt Date) After(u Date) bool {
	return dt.src.After(u.src)
}

// Before behavior as time.Time
func (dt Date) Before(u Date) bool {

	return dt.src.Before(u.src)
}

// Equal behavior as time.Time
func (dt Date) Equal(u Date) bool {
	return dt.src.Equal(u.src)
}

// Time convert to time.Time
func (dt Date) Time() time.Time {
	return dt.src
}

// AddDate behavior as time.Time
func (dt Date) AddDate(years int, months int, days int) Date {
	return NewDateFromTime(dt.src.AddDate(years, months, days))
}

// Sub  behavior as time.Time
func (dt Date) Sub(u Date) time.Duration {
	return dt.src.Sub(u.src)
}

// Add behavior as time.Time
func (dt Date) Add(d time.Duration) Date {
	return NewDateFromTime(dt.src.Add(d))
}

// IsZero behavior as time.Time
func (dt Date) IsZero() bool {
	return dt.src.IsZero()
}

// Date  behavior as time.Time
func (dt Date) Date() (year int, month time.Month, day int) {
	return dt.src.Date()
}

// Year behavior as time.Time
func (dt Date) Year() int {
	return dt.src.Year()
}

// Month behavior as time.Time
func (dt Date) Month() time.Month {
	return dt.src.Month()
}

// Day behavior as time.Time
func (dt Date) Day() int {
	return dt.src.Day()
}

// Weekday behavior as time.Time
func (dt Date) Weekday() time.Weekday {
	return dt.src.Weekday()
}

// ISOWeek behavior as time.Time
func (dt Date) ISOWeek() (year int, week int) {
	return dt.src.ISOWeek()
}

// YearDay behavior as time.Time
func (dt Date) YearDay() int {
	return dt.src.YearDay()
}

// UnmarshalText behavior as time.Time
func (dt *Date) UnmarshalText(text []byte) error {
	t, err := time.Parse(dateFormatLayout, string(text))
	if err != nil {
		return err
	}
	dt.src = t
	return nil
}

// MarshalText behavior as time.Time
func (dt Date) MarshalText() ([]byte, error) {
	return []byte(dt.src.Format(dateFormatLayout)), nil
}

const dateFormatLayout = "2006-01-02"

var _ driver.Valuer = Date{}
var _ sql.Scanner = &Date{}
var _ encoding.TextUnmarshaler = &Date{}

// Scan for sql.Scanner
func (dt *Date) Scan(value interface{}) error {
	src, ok := value.(time.Time)

	var dst Date
	if ok {

		dst = NewDateFromTime(src)

	} else {
		src, ok := value.([]byte)
		if !ok {
			return ErrInvalidValueType
		}
		t, err := time.Parse(dateFormatLayout, string(src))
		if err != nil {
			return err
		}
		dst = NewDateFromTime(t)
	}
	dt.src = dst.src
	return nil
}

// Value for driver.Valuer
func (dt Date) Value() (driver.Value, error) {
	return dt.src, nil
}
