package mysqltype

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"time"
)

// DateTime support MySQL DateTime type
// https://dev.mysql.com/doc/refman/8.0/en/datetime.html
type DateTime struct {
	src time.Time `gorm:"type:datetime"`
}

// NewDateTime Create new DateTime from time.Date
func NewDateTime(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime {
	t := time.Date(year, month, day, hour, min, sec, nsec, loc)
	return NewDateTimeFromTime(t)
}

// NewDateTimeFromTime Create new DateTime from Time
func NewDateTimeFromTime(t time.Time) DateTime {
	return DateTime{src: t}
}

// MinDateTime Minimum DateTime
func MinDateTime() DateTime {
	return NewDateTimeFromTime(time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// MaxDateTime Maximum DateTime
func MaxDateTime() DateTime {
	return NewDateTimeFromTime(time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC))
}

// NowDateTime Create Now time for MySQL DataBase
func NowDateTime() DateTime {
	return NewDateTimeFromTime(time.Now())
}

// After behavior as time.Time
func (dt DateTime) After(u DateTime) bool {
	return dt.src.After(u.src)
}

// Before behavior as time.Time
func (dt DateTime) Before(u DateTime) bool {

	return dt.src.Before(u.src)
}

// Equal behavior as time.Time
func (dt DateTime) Equal(u DateTime) bool {
	return dt.src.Equal(u.src)
}

// Time convert to time.Time
func (dt DateTime) Time() time.Time {
	return dt.src
}

// Round behavior as time.Time
func (dt DateTime) Round(d time.Duration) DateTime {
	return NewDateTimeFromTime(dt.src.Round(d))
}

// UnixNano behavior as time.Time
func (dt DateTime) UnixNano() int64 {
	return dt.src.UnixNano()
}

// Unix behavior as time.Time
func (dt DateTime) Unix() int64 {
	return dt.src.Unix()
}

// AddDate behavior as time.Time
func (dt DateTime) AddDate(years int, months int, days int) DateTime {
	return NewDateTimeFromTime(dt.src.AddDate(years, months, days))
}

// Sub  behavior as time.Time
func (dt DateTime) Sub(u DateTime) time.Duration {
	return dt.src.Sub(u.src)
}

// Add behavior as time.Time
func (dt DateTime) Add(d time.Duration) DateTime {
	return NewDateTimeFromTime(dt.src.Add(d))
}

// Location behavior as time.Time
func (dt DateTime) Location() *time.Location {
	return dt.src.Location()
}

// Local behavior as time.Time
func (dt DateTime) Local() DateTime {
	return NewDateTimeFromTime(dt.src.Local())
}

// UTC behavior as time.Time
func (dt DateTime) UTC() DateTime {
	return NewDateTimeFromTime(dt.src.UTC())
}

// In behavior as time.Time
func (dt DateTime) In(loc *time.Location) DateTime {
	return NewDateTimeFromTime(dt.src.In(loc))
}

// IsZero behavior as time.Time
func (dt DateTime) IsZero() bool {
	return dt.src.IsZero()
}

// Date  behavior as time.Time
func (dt DateTime) Date() (year int, month time.Month, day int) {
	return dt.src.Date()
}

// Year behavior as time.Time
func (dt DateTime) Year() int {
	return dt.src.Year()
}

// Month behavior as time.Time
func (dt DateTime) Month() time.Month {
	return dt.src.Month()
}

// Day behavior as time.Time
func (dt DateTime) Day() int {
	return dt.src.Day()
}

// Weekday behavior as time.Time
func (dt DateTime) Weekday() time.Weekday {
	return dt.src.Weekday()
}

// ISOWeek behavior as time.Time
func (dt DateTime) ISOWeek() (year int, week int) {
	return dt.src.ISOWeek()
}

// Clock behavior as time.Time
func (dt DateTime) Clock() (hour, min, sec int) {
	return dt.src.Clock()
}

// Hour behavior as time.Time
func (dt DateTime) Hour() int {
	return dt.src.Hour()
}

// Minute behavior as time.Time
func (dt DateTime) Minute() int {
	return dt.src.Minute()
}

// Second behavior as time.Time
func (dt DateTime) Second() int {
	return dt.src.Second()
}

// Nanosecond behavior as time.Time
func (dt DateTime) Nanosecond() int {
	return dt.src.Nanosecond()
}

// Truncate behavior as time.Time
func (dt DateTime) Truncate(d time.Duration) DateTime {
	return NewDateTimeFromTime(dt.src.Truncate(d))
}

// YearDay behavior as time.Time
func (dt DateTime) YearDay() int {
	return dt.src.YearDay()
}

// UnmarshalText behavior as time.Time
func (dt *DateTime) UnmarshalText(text []byte) error {
	return dt.src.UnmarshalText(text)
}

// MarshalText behavior as time.Time
func (dt DateTime) MarshalText() ([]byte, error) {
	return dt.src.MarshalText()
}

const dateTimeFormatLayout = "2006-01-02 15:04:05.999999999"

var _ driver.Valuer = DateTime{}
var _ sql.Scanner = &DateTime{}
var _ encoding.TextUnmarshaler = &DateTime{}

// Scan for sql.Scanner
func (dt *DateTime) Scan(value interface{}) error {
	src, ok := value.(time.Time)

	var dst DateTime
	if ok {

		dst = NewDateTimeFromTime(src)

	} else {
		src, ok := value.([]byte)
		if !ok {
			return ErrInvalidValueType
		}
		t, err := time.Parse(dateTimeFormatLayout, string(src))
		if err != nil {
			return err
		}
		dst = NewDateTimeFromTime(t)
	}
	dt.src = dst.src
	return nil
}

// Value for driver.Valuer
func (dt DateTime) Value() (driver.Value, error) {
	return dt.src, nil
}
