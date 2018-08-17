package mysqltype

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DateTimeFieldTestStruct struct {
	ID         int
	TargetDate DateTime `gorm:"not null"`
}

type DateTimeNullFieldTestStruct struct {
	ID         int
	TargetDate *DateTime
}

func TestDateTimeField(t *testing.T) {
	t.Parallel()
	dateTime := NowDateTime().Truncate(1 * time.Second)
	target := &DateTimeFieldTestStruct{
		TargetDate: dateTime,
	}
	assert.NoError(t, DB.AutoMigrate(target).Error)
	assert.NoError(t, DB.Create(target).Error)

	dummy := &DateTimeNullFieldTestStruct{}

	assertMySQLErrNumber(t, DB.Table("date_time_field_test_structs").Create(&dummy).Error, mySQLNullError)
	dst := &DateTimeFieldTestStruct{
		ID: target.ID,
	}
	assert.NoError(t, DB.First(dst).Error)
	assertTimeEquals(t, target.TargetDate.src, dst.TargetDate.src)

	assert.NoError(t, DB.Save(&DateTimeFieldTestStruct{TargetDate: MaxDateTime()}).Error)
	assert.NoError(t, DB.Save(&DateTimeFieldTestStruct{TargetDate: MinDateTime()}).Error)
}

func TestDateTimeFieldLocale(t *testing.T) {
	assert.NoError(t, DB.AutoMigrate(DateTimeFieldTestStruct{}).Error)
	asiaTokyo, err := time.LoadLocation("Asia/Tokyo")
	assert.NoError(t, err)
	nowUTC := NowDateTime().UTC()
	nowJST := nowUTC.In(asiaTokyo)

	vUTC := DateTimeFieldTestStruct{TargetDate: nowUTC}
	vJST := DateTimeFieldTestStruct{TargetDate: nowJST}

	assert.NoError(t, DB.Save(&vUTC).Error)
	assert.NoError(t, DB.Save(&vJST).Error)
	dst1 := DateTimeFieldTestStruct{ID: vUTC.ID}
	dst2 := DateTimeFieldTestStruct{ID: vJST.ID}

	assert.NoError(t, DB.Find(&dst1).Error)
	assert.NoError(t, DB.Find(&dst2).Error)
	assert.Equal(t, dst1.TargetDate, dst2.TargetDate)
}

func TestDateTimeMarshalJSON(t *testing.T) {
	now := NowDateTime()
	expected, err := now.MarshalText()
	assert.NoError(t, err)
	actual, err := json.Marshal(now)
	assert.NoError(t, err)
	assert.EqualValues(t, "\""+string(expected)+"\"", string(actual))
}

func TestDateTimeValue(t *testing.T) {
	t.Parallel()
	dateTime := NowDateTime()
	v, err := dateTime.Value()
	assert.NoError(t, err)
	assert.Equal(t, dateTime.src, v)
}

func TestDateTimeScan(t *testing.T) {
	t.Parallel()
	target := DateTime{}
	now := time.Now()
	assert.NoError(t, target.Scan(now))
	assertTimeEquals(t, now, target.src)
	nowStr := now.Format(dateTimeFormatLayout)
	nowFromFormat, err := time.Parse(dateTimeFormatLayout, nowStr)
	assert.NoError(t, err)
	target2 := DateTime{}
	assert.NoError(t, target2.Scan([]byte(nowStr)))
	assertTimeEquals(t, nowFromFormat, target2.src)
}

func TestDateTimeAfterAndBefore(t *testing.T) {
	t.Parallel()
	v1 := NewDateTime(2008, 10, 12, 2, 30, 2, 0, time.UTC)
	v2 := NewDateTime(2008, 10, 12, 2, 30, 3, 0, time.UTC)

	assert.True(t, v2.After(v1))
	assert.False(t, v1.After(v2))
	assert.True(t, v1.Before(v2))
	assert.False(t, v2.Before(v1))
}

func TestDateTimeEqual(t *testing.T) {
	t.Parallel()
	assert.True(t, MinDateTime().Equal(MinDateTime()))
	assert.False(t, NowDateTime().Equal(MinDateTime()))
	assert.False(t, MinDateTime().Equal(NowDateTime()))
}

func TestDateTimeToTime(t *testing.T) {
	t.Parallel()
	min := MinDateTime().Time()
	assertTimeEquals(t, MinDateTime().src, min)
	max := MaxDateTime().Time()
	assertTimeEquals(t, MaxDateTime().src, max)
}

func TestDateTimeTextUnmarshalText(t *testing.T) {
	t.Parallel()
	text := "2016-12-31T20:02:05.123456Z"
	target := NewDateTimeFromTime(time.Time{})
	target.UnmarshalText([]byte(text))
	expected := time.Time{}
	expected.UnmarshalText([]byte(text))
	assertTimeEquals(t, expected, target.src)
}

func TestDateTimeUnixNano(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	assert.Equal(t, now.src.UnixNano(), now.UnixNano())
}

func TestDateTimeUnix(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	assert.Equal(t, now.src.Unix(), now.Unix())
}

func TestDateTimeRound(t *testing.T) {
	t.Parallel()
	v := NewDateTime(2008, 10, 12, 10, 3, 8, 8884, time.UTC)
	expected := NewDateTime(2008, 10, 12, 10, 3, 10, 0, time.UTC)
	rounded := v.Round(5 * time.Second)
	assertTimeEquals(t, expected.src, rounded.src)
}

func TestDateTimeTrauncate(t *testing.T) {
	t.Parallel()
	v := NewDateTime(2008, 10, 12, 10, 3, 9, 8884, time.UTC)
	expected := NewDateTime(2008, 10, 12, 10, 3, 5, 0, time.UTC)
	truncated := v.Truncate(5 * time.Second)
	assertTimeEquals(t, expected.src, truncated.src)
}

func TestDateTimeAddDate(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	expected := now.src.AddDate(1, 2, 4)
	actual := now.AddDate(1, 2, 4).src
	assertTimeEquals(t, expected, actual)
}

func TestDateTimeSub(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	sub := NewDateTime(0, 0, 0, 10, 0, 0, 0, time.UTC)
	expected := now.src.Sub(sub.src)
	actual := now.Sub(sub)
	assert.Equal(t, expected, actual)
}

func TestDateTimeAdd(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	add := 12 * time.Second
	expected := now.src.Add(add)
	actual := now.Add(add)
	assert.Equal(t, expected, actual.src)
}

func TestDateTimeLocation(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	assert.Equal(t, now.src.Location(), now.Location())
}

func TestDateTimeLocal(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	assert.Equal(t, now.Local().src, now.src.Local())
}

func TestDateTimeIn(t *testing.T) {
	t.Parallel()
	asiaTokyo, err := time.LoadLocation("Asia/Tokyo")
	assert.NoError(t, err)
	now := NewDateTimeFromTime(time.Now().UTC())
	assert.Equal(t, now.In(asiaTokyo).src, now.src.In(asiaTokyo))
}

func TestDateTimeUTC(t *testing.T) {
	t.Parallel()
	asiaTokyo, err := time.LoadLocation("Asia/Tokyo")
	assert.NoError(t, err)
	nowJST := NewDateTimeFromTime(time.Now().In(asiaTokyo))
	assert.Equal(t, nowJST.UTC().src, nowJST.src.UTC())
}

func TestDateTimeIsZero(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	zero := DateTime{}

	assert.False(t, now.IsZero())
	assert.True(t, zero.IsZero())
}

func TestDateTimeDate(t *testing.T) {
	t.Parallel()
	now := NowDateTime()
	expectedYear, expectedMonth, expectedDay := now.src.Date()
	actualYear, actualMonth, actualDay := now.Date()
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, expectedMonth, actualMonth)
	assert.Equal(t, expectedDay, actualDay)
}

func TestDateTimeYear(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	assert.Equal(t, max.src.Year(), max.Year())
	assert.Equal(t, 9999, max.Year())
}

func TestDateTimeMonth(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	assert.Equal(t, max.src.Month(), max.Month())
	assert.Equal(t, time.Month(12), max.Month())
}

func TestDateTimeDay(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	assert.Equal(t, max.src.Day(), max.Day())
	assert.Equal(t, 31, max.Day())
}

func TestDateTimeWeekDay(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	assert.Equal(t, max.src.Weekday(), max.Weekday())
	assert.Equal(t, time.Friday, max.Weekday())
}

func TestDateTimeISOWeek(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	expectedYear, expectedWeek := max.src.ISOWeek()
	actualYear, actualWeek := max.ISOWeek()
	assert.Equal(t, expectedYear, actualYear)
	assert.Equal(t, expectedWeek, actualWeek)
	assert.Equal(t, 9999, actualYear)
	assert.Equal(t, 52, actualWeek)
}

func TestDateTimeClock(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	expectedHour, expectedMin, expectedSec := max.src.Clock()
	actualHour, actualMin, actualSec := max.Clock()
	assert.Equal(t, expectedHour, actualHour)
	assert.Equal(t, expectedMin, actualMin)
	assert.Equal(t, expectedSec, actualSec)
	assert.Equal(t, 23, actualHour)
	assert.Equal(t, 59, actualMin)
	assert.Equal(t, 59, actualSec)
}

func TestDateTimeHour(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	expected := max.src.Hour()
	actual := max.Hour()
	assert.Equal(t, expected, actual)
	assert.Equal(t, 23, actual)
}

func TestDateTimeMinute(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	expected := max.src.Minute()
	actual := max.Minute()
	assert.Equal(t, expected, actual)
	assert.Equal(t, 59, actual)
}

func TestDateTimeSecond(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	expected := max.src.Second()
	actual := max.Second()
	assert.Equal(t, expected, actual)
	assert.Equal(t, 59, actual)
}

func TestDateTimeNanosecond(t *testing.T) {
	t.Parallel()
	max := MaxDateTime()
	expected := max.src.Nanosecond()
	actual := max.Nanosecond()
	assert.Equal(t, expected, actual)
	assert.Equal(t, 999999, actual)
}
