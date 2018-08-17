package mysqltype

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DateFieldTestStruct struct {
	ID         int
	TargetDate Date `gorm:"not null"`
}

type DateNullFieldTestStruct struct {
	ID         int
	TargetDate *Date
}

func TestDateField(t *testing.T) {
	t.Parallel()
	dateTime := NowDate()
	target := &DateFieldTestStruct{
		TargetDate: dateTime,
	}
	assert.NoError(t, DB.AutoMigrate(target).Error)
	assert.NoError(t, DB.Create(target).Error)

	dummy := &DateNullFieldTestStruct{}

	assertMySQLErrNumber(t, DB.Table("date_field_test_structs").Create(&dummy).Error, mySQLNullError)
	dst := &DateFieldTestStruct{
		ID: target.ID,
	}
	assert.NoError(t, DB.First(dst).Error)
	assertTimeEquals(t, target.TargetDate.src, dst.TargetDate.src)

	assert.NoError(t, DB.Save(&DateFieldTestStruct{TargetDate: MaxDate()}).Error)
	assert.NoError(t, DB.Save(&DateFieldTestStruct{TargetDate: MinDate()}).Error)
}

func TestDateMarshalJSON(t *testing.T) {
	now := NowDate()
	expected, err := now.MarshalText()
	assert.NoError(t, err)
	actual, err := json.Marshal(now)
	assert.NoError(t, err)
	assert.EqualValues(t, "\""+string(expected)+"\"", string(actual))
}

func TestDateValue(t *testing.T) {
	t.Parallel()
	dateTime := NowDate()
	v, err := dateTime.Value()
	assert.NoError(t, err)
	assert.Equal(t, dateTime.src, v)
}

func TestDateScan(t *testing.T) {
	t.Parallel()
	target := Date{}
	now := NowDate().src
	assert.NoError(t, target.Scan(now))
	assertTimeEquals(t, now, target.src)
	nowStr := now.Format(dateFormat)
	nowFromFormat, err := time.Parse(dateFormat, nowStr)
	assert.NoError(t, err)
	target2 := Date{}
	assert.NoError(t, target2.Scan([]byte(nowStr)))
	assertTimeEquals(t, nowFromFormat, target2.src)
}

func TestDateAfterAndBefore(t *testing.T) {
	t.Parallel()
	v1 := NewDate(2008, 10, 12)
	v2 := NewDate(2008, 10, 13)

	assert.True(t, v2.After(v1))
	assert.False(t, v1.After(v2))
	assert.True(t, v1.Before(v2))
	assert.False(t, v2.Before(v1))
}

func TestDateEqual(t *testing.T) {
	t.Parallel()
	assert.True(t, MinDate().Equal(MinDate()))
	assert.False(t, NowDate().Equal(MinDate()))
	assert.False(t, MinDate().Equal(NowDate()))
}

func TestDateToTime(t *testing.T) {
	t.Parallel()
	min := MinDate().Time()
	assertTimeEquals(t, MinDate().src, min)
	max := MaxDate().Time()
	assertTimeEquals(t, MaxDate().src, max)
}

func TestDateTextUnmarshalText(t *testing.T) {
	t.Parallel()
	text := "2016-12-31T20:02:05.123456Z"
	target := NewDateFromTime(time.Time{})
	target.UnmarshalText([]byte(text))
	expected := time.Time{}
	expected.UnmarshalText([]byte(text))
	assertTimeEquals(t, expected, target.src)
}

func TestDateAddDate(t *testing.T) {
	t.Parallel()
	now := NowDate()
	expected := now.src.AddDate(1, 2, 4)
	actual := now.AddDate(1, 2, 4).src
	assertTimeEquals(t, expected, actual)
}

func TestDateSub(t *testing.T) {
	t.Parallel()
	now := NowDate()
	sub := NewDate(0, 0, 1)
	expected := now.src.Sub(sub.src)
	actual := now.Sub(sub)
	assert.Equal(t, expected, actual)
}

func TestDateAdd(t *testing.T) {
	t.Parallel()
	now := NowDate()
	add := 24 * time.Hour * 2
	expected := now.src.Add(add)
	actual := now.Add(add)
	assert.Equal(t, expected, actual.src)
}

func TestDateIsZero(t *testing.T) {
	t.Parallel()
	now := NowDate()
	zero := Date{}

	assert.False(t, now.IsZero())
	assert.True(t, zero.IsZero())
}

func TestDateDate(t *testing.T) {
	t.Parallel()
	now := NowDate()
	expectedYear, expectedMonth, expectedDay := now.src.Date()
	actualYear, actualMonth, actualDay := now.Date()
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, expectedMonth, actualMonth)
	assert.Equal(t, expectedDay, actualDay)
}

func TestDateYear(t *testing.T) {
	t.Parallel()
	max := MaxDate()
	assert.Equal(t, max.src.Year(), max.Year())
	assert.Equal(t, 9999, max.Year())
}

func TestDateMonth(t *testing.T) {
	t.Parallel()
	max := MaxDate()
	assert.Equal(t, max.src.Month(), max.Month())
	assert.Equal(t, time.Month(12), max.Month())
}

func TestDateDay(t *testing.T) {
	t.Parallel()
	max := MaxDate()
	assert.Equal(t, max.src.Day(), max.Day())
	assert.Equal(t, 31, max.Day())
}

func TestDateWeekDay(t *testing.T) {
	t.Parallel()
	max := MaxDate()
	assert.Equal(t, max.src.Weekday(), max.Weekday())
	assert.Equal(t, time.Friday, max.Weekday())
}

func TestDateISOWeek(t *testing.T) {
	t.Parallel()
	max := MaxDate()
	expectedYear, expectedWeek := max.src.ISOWeek()
	actualYear, actualWeek := max.ISOWeek()
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, expectedWeek, actualWeek)
	assert.Equal(t, 9999, actualYear)
	assert.Equal(t, 52, actualWeek)
}
