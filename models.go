package bublyk

import (
	"errors"
	"time"

	faststrconv "github.com/kaatinga/strconv"
)

const (
	monthMask = 0b0000000111100000
	dayMask   = 0b0000000000011111

	maximumDate Date = 0b1111111110011111 // 2127-12-31
	minimumDate Date = 0b0000000000100001 // 2000-01-01
	noDate      Date = 0

	postgreSQLFormat = "2006-01-02"
)

func Now() Date {
	now := time.Now().UTC()
	return NewDateFromTime(&now)
}

func CurrentMonth() Date {
	now := Now()
	return NewDate(now.Year(), now.Month(), 1)
}

// Date represents a calendar date starting 2000 year and finishing the year 2127.
type Date uint16

func (thisDate Date) Year() uint16 {
	return uint16(thisDate>>9) + 2000
}

func (thisDate Date) Month() byte {
	return byte((thisDate & monthMask) >> 5)
}

func (thisDate Date) Day() byte {
	return byte(thisDate & dayMask)
}

func (thisDate Date) IsSet() bool {
	return thisDate != 0
}

func (thisDate Date) IsFuture() bool {
	return thisDate > Now()
}

// MonthAfter checks whether the date at least one month after the target date or not.
func (thisDate Date) MonthAfter(date Date) bool {
	if thisDate.Year() == date.Year() {
		return thisDate.Month() > date.Month()
	}
	return thisDate.Year() > date.Year()
}

// MonthBefore checks whether the date at least one month before the target date or not.
func (thisDate Date) MonthBefore(date Date) bool {
	if thisDate.Year() == date.Year() {
		return thisDate.Month() < date.Month()
	}
	return thisDate.Year() < date.Year()
}

// String returns date as string in the default PostgreSQL date format, YYYY-MM-DD.
func (thisDate Date) String() string {
	if !thisDate.IsSet() {
		return "null"
	}
	return format[string](thisDate)
}

// format prepares a binary or text of the date in PostgreSQL format.
func format[V []byte | string](thisDate Date) V {
	month, day, year := getDateAsBinaries(thisDate)
	var output = make([]byte, 10)
	output[0] = year[0]
	output[1] = year[1]
	output[2] = year[2]
	output[3] = year[3]
	output[4] = '-'

	switch len(month) {
	case 2:
		output[5] = month[0]
		output[6] = month[1]
	default:
		output[5] = 48
		output[6] = month[0]
	}
	output[7] = '-'

	switch len(day) {
	case 2:
		output[8] = day[0]
		output[9] = day[1]
	default:
		output[8] = 48
		output[9] = day[0]
	}
	return V(output)
}

func getDateAsBinaries(thisDate Date) ([]byte, []byte, []byte) {
	var month = faststrconv.Byte2Bytes(thisDate.Month())
	var day = faststrconv.Byte2Bytes(thisDate.Day())
	var year = faststrconv.Uint162Bytes(thisDate.Year())
	return month, day, year
}

// DMYWithDots returns date as string in the DD.MM.YYYY format.
func (thisDate Date) DMYWithDots() string {
	month, day, year := getDateAsBinaries(thisDate)
	var output = make([]byte, 10)

	switch len(day) {
	case 2:
		output[0] = day[0]
		output[1] = day[1]
	default:
		output[0] = 48
		output[1] = day[0]
	}
	output[2] = '.'

	switch len(month) {
	case 2:
		output[3] = month[0]
		output[4] = month[1]
	default:
		output[3] = 48
		output[4] = month[0]
	}
	output[5] = '.'

	output[6] = year[0]
	output[7] = year[1]
	output[8] = year[2]
	output[9] = year[3]

	return string(output)
}

func (thisDate Date) Format(layout string) string {
	switch layout {
	case postgreSQLFormat:
		return thisDate.String()
	default:
		return makeTime(thisDate.Year(), thisDate.Month(), thisDate.Day()).Format(layout)
	}
}

func Parse(formattedDate string) (Date, error) {
	if len([]rune(formattedDate)) != len([]rune(postgreSQLFormat)) {
		return 0, errors.New("incorrect date format")
	}

	year, err := faststrconv.GetUint16(formattedDate[0:4])
	if err != nil {
		return 0, err
	}

	var month byte
	month, err = faststrconv.GetByte(formattedDate[5:7])
	if err != nil {
		return 0, err
	}

	var day byte
	day, err = faststrconv.GetByte(formattedDate[8:10])
	if err != nil {
		return 0, err
	}

	return NewDate(year, month, day), nil
}

func (thisDate Date) Time() *time.Time {
	return makeTime(thisDate.Year(), thisDate.Month(), thisDate.Day())
}

func (thisDate Date) NextDay() Date {
	if thisDate.Day() > 27 {
		timeDate := makeTime(thisDate.Year(), thisDate.Month(), thisDate.Day()).AddDate(0, 0, 1)
		return NewDateFromTime(&timeDate)
	}
	return thisDate&^dayMask | (thisDate&dayMask + 1)
}

func (thisDate Date) PreviousDay() Date {
	if thisDate.Day() == 1 {
		timeDate := makeTime(thisDate.Year(), thisDate.Month(), thisDate.Day()).AddDate(0, 0, -1)
		return NewDateFromTime(&timeDate)
	}
	return thisDate&^dayMask | (thisDate&dayMask - 1)
}

func (thisDate Date) NextWeek() Date {
	if thisDate.Day() > 21 {
		timeDate := makeTime(thisDate.Year(), thisDate.Month(), thisDate.Day()).AddDate(0, 0, 7)
		return NewDateFromTime(&timeDate)
	}
	return thisDate&^dayMask | (thisDate&dayMask + 7)
}

func (thisDate Date) PreviousWeek() Date {
	if thisDate.Day() < 8 {
		timeDate := makeTime(thisDate.Year(), thisDate.Month(), thisDate.Day()).AddDate(0, 0, -7)
		return NewDateFromTime(&timeDate)
	}
	return thisDate&^dayMask | (thisDate&dayMask - 7)
}

func NewDate(year uint16, month, day byte) Date {
	if year < 2000 {
		return minimumDate
	}
	if year > 2127 {
		return maximumDate
	}
	if day > 28 || month > 12 || day == 0 || month == 0 {
		yearInt, monthMonth, dayInt := makeTime(year, month, day).Date()
		year, month, day = uint16(yearInt), byte(monthMonth), byte(dayInt)
	}
	return composeDate(year, month, day)
}

func makeTime(year uint16, month, day byte) *time.Time {
	newTime := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	return &newTime
}

// NewDateFromTime create new date using time.Time model.
func NewDateFromTime(t *time.Time) Date {
	if t == nil {
		return noDate
	}
	year := t.Year()
	if year < 2000 {
		return minimumDate
	}
	if year > 2127 {
		return maximumDate
	}
	return composeDate(uint16(year), byte(t.Month()), byte(t.Day()))
}

func composeDate(year uint16, month, day byte) Date {
	return (Date(year-2000) << 9) + (Date(month) << 5) + Date(day)
}
