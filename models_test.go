package bublyk

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	tests := []struct {
		year, month, day int
		want             string
	}{
		{2021, 2, 28, "2021-02-28"},
		{2001, 2, 28, "2001-02-28"},
		{2021, 12, 31, "2021-12-31"},
		{2021, 1, 1, "2021-01-01"},
		{2000, 1, 1, "2000-01-01"},
		{2127, 12, 31, "2127-12-31"},
		{2031, 1, 1, "2031-01-01"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%04d-%02d-%02d", tt.year, tt.month, tt.day), func(t *testing.T) {
			newTime := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
			date := NewDateFromTime(&newTime)
			if uint16(tt.year) != date.Year() {
				t.Errorf("Year is incorrect.\nhave %v\nwant %v", date.Year(), tt.year)
			}

			if byte(tt.month) != date.Month() {
				t.Log("error NewDateFromTime", date.Month())
				t.Errorf("Month is incorrect.\nhave %v\nwant %v", date.Month(), tt.month)
			}

			if byte(tt.day) != date.Day() {
				t.Errorf("Day is incorrect.\nhave %v\nwant %v", date.Day(), tt.day)
			}

			date2 := NewDate(uint16(tt.year), byte(tt.month), byte(tt.day))
			if uint16(tt.year) != date2.Year() {
				t.Errorf("Year is incorrect.\nhave %v\nwant %v", date.Year(), tt.year)
			}

			if byte(tt.month) != date2.Month() {
				t.Logf("error NewDate")
				t.Errorf("Month is incorrect.\nhave %v\nwant %v", date.Month(), tt.month)
			}

			if byte(tt.day) != date2.Day() {
				t.Errorf("Day is incorrect.\nhave %v\nwant %v", date.Day(), tt.day)
			}

			if tt.want != date.String() {
				t.Errorf("String() = %v, want %v", date.String(), tt.want)
			}

			t.Logf("%16b\n", date)
			t.Log(date.String())
		})
	}

	// Test for the maximum date.
	d := NewDate(2227, 12, 31)
	if maximumDate != d {
		t.Errorf("maximumDate = %v, want %v", d, maximumDate)
	}

	// Test for the minimum date.
	d = NewDate(1999, 1, 1)
	if minimumDate != d {
		t.Errorf("minimumDate = %v, want %v", d, minimumDate)
	}
}

func TestDate_Format(t *testing.T) {
	tests := []struct {
		this   Date
		layout string
		want   string
	}{
		{maximumDate, postgreSQLFormat, "2127-12-31"},
		{maximumDate, time.RFC822, "31 Dec 27 00:00 UTC"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.this.Format(tt.layout); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextDay(t *testing.T) {
	testCases := []struct {
		this, want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2022, 1, 1)},
		{NewDate(2021, 12, 1), NewDate(2021, 12, 2)},
		{NewDate(2021, 2, 28), NewDate(2021, 3, 1)},
	}
	for _, tt := range testCases {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.NextDay(); got != tt.want {
				t.Errorf("NextDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousDay(t *testing.T) {
	var testCases = []struct {
		this Date
		want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2021, 12, 30)},
		{NewDate(2021, 12, 2), NewDate(2021, 12, 1)},
		{NewDate(2021, 3, 1), NewDate(2021, 2, 28)},
		{NewDate(2021, 3, 2), NewDate(2021, 3, 1)},
		{minimumDate, minimumDate},
	}
	for _, tt := range testCases {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.PreviousDay(); got != tt.want {
				t.Errorf("PreviousDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextWeek(t *testing.T) {
	tests := []struct {
		this, want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2022, 1, 7)},
		{NewDate(2021, 12, 1), NewDate(2021, 12, 8)},
		{NewDate(2021, 2, 28), NewDate(2021, 3, 7)},
		{NewDate(2021, 2, 1), NewDate(2021, 2, 8)},
		{NewDate(2127, 12, 31), maximumDate},
		{maximumDate, maximumDate},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.NextWeek(); got != tt.want {
				t.Errorf("NextWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousWeek(t *testing.T) {
	tests := []struct {
		this, want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2021, 12, 24)},
		{NewDate(2022, 1, 3), NewDate(2021, 12, 27)},
		{NewDate(2021, 12, 1), NewDate(2021, 11, 24)},
		{NewDate(2021, 3, 2), NewDate(2021, 2, 23)},
		{NewDate(2021, 3, 9), NewDate(2021, 3, 2)},
		{NewDate(2000, 1, 1), minimumDate},
		{minimumDate, minimumDate},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.PreviousWeek(); got != tt.want {
				t.Errorf("PreviousWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextMonth(t *testing.T) {
	tests := []struct {
		date, want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2022, 1, 31)},
		{NewDate(2024, 2, 29), NewDate(2024, 3, 29)},
		{NewDate(2024, 1, 31), NewDate(2024, 3, 02)},
		{NewDate(2021, 12, 1), NewDate(2022, 1, 1)},
		{NewDate(2021, 3, 2), NewDate(2021, 4, 2)},
		{NewDate(2023, 1, 1), NewDate(2023, 2, 1)},
		{NewDate(2127, 12, 1), maximumDate},
		{maximumDate, maximumDate},
	}
	for _, tt := range tests {
		t.Run(tt.date.String(), func(t *testing.T) {
			if got := tt.date.NextMonth(); got != tt.want {
				t.Errorf("NextMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousMonth(t *testing.T) {
	tests := []struct {
		this Date
		want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2021, 12, 1)},
		{NewDate(2024, 2, 29), NewDate(2024, 1, 29)},
		{NewDate(2021, 12, 1), NewDate(2021, 11, 1)},
		{NewDate(2021, 3, 2), NewDate(2021, 2, 2)},
		{NewDate(2024, 1, 1), NewDate(2023, 12, 1)},
		{NewDate(1999, 1, 1), minimumDate},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.PreviousMonth(); got != tt.want {
				t.Errorf("PreviousMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		formattedDate string
		want          Date
		wantErr       bool
	}{
		{"2021-12-01", NewDate(2021, 12, 1), false},
	}
	for _, tt := range tests {
		t.Run(tt.formattedDate, func(t *testing.T) {
			got, err := Parse(tt.formattedDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsFuture(t *testing.T) {
	tests := []struct {
		date Date
		want bool
	}{
		{NewDate(2020, 12, 12), false},
		{NewDate(2035, 12, 12), true},
	}
	for _, tt := range tests {
		t.Run(tt.date.String(), func(t *testing.T) {
			if got := tt.date.IsFuture(); got != tt.want {
				t.Errorf("IsFuture() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_binary(t *testing.T) {
	tests := []struct {
		thisDate Date
		want     string
	}{
		{NewDate(2022, 10, 10), "2022-10-10"},
		{NewDate(2000, 1, 1), "2000-01-01"},
	}
	for _, tt := range tests {
		t.Run(tt.thisDate.String(), func(t *testing.T) {
			receivedDate := format[string](tt.thisDate)
			if receivedDate != tt.want {
				t.Errorf("format() = %v, want %v", receivedDate, tt.want)
			}
		})
	}
}

func TestCurrentMonth(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"now1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurrentMonth()
			if got.Day() != 1 {
				t.Errorf("CurrentMonth() returned incorrect day = %v", got.Day())
			}
		})
	}
}

func TestDate_DMYWithDots(t *testing.T) {
	tests := []struct {
		thisDate Date
		want     string
	}{
		{NewDate(2000, 1, 1), "01.01.2000"},
		{NewDate(2022, 11, 11), "11.11.2022"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.thisDate.DMYWithDots(); got != tt.want {
				t.Errorf("DMYWithDots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		thisDate Date
		want     string
	}{
		{0, "null"},
		{NewDate(2022, 10, 10), "2022-10-10"},
	}
	for _, tt := range tests {
		t.Run(tt.thisDate.String(), func(t *testing.T) {
			if got := tt.thisDate.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
