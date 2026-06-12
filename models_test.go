package bublyk

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	cases := []struct {
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
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%04d-%02d-%02d", tt.year, tt.month, tt.day), func(t *testing.T) {
			newTime := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
			date := NewDateFromTime(newTime)
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
	cases := []struct {
		this   Date
		layout string
		want   string
	}{
		{maximumDate, postgreSQLFormat, "2127-12-31"},
		{maximumDate, time.RFC822, "31 Dec 27 00:00 UTC"},
		{Date(0), postgreSQLFormat, "null"},
		{Date(0), time.RFC822, "null"},
	}
	for _, tt := range cases {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.this.Format(tt.layout); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextDay(t *testing.T) {
	cases := []struct {
		name string
		this Date
		want Date
	}{
		{"Regular day end of month", NewDate(2021, 12, 31), NewDate(2022, 1, 1)},
		{"Regular day middle of month", NewDate(2021, 12, 1), NewDate(2021, 12, 2)},
		{"End of February non-leap year", NewDate(2021, 2, 28), NewDate(2021, 3, 1)},
		{"End of February leap year", NewDate(2024, 2, 29), NewDate(2024, 3, 1)},
		{"End of 30-day month", NewDate(2021, 4, 30), NewDate(2021, 5, 1)},
		{"End of 31-day month", NewDate(2021, 7, 31), NewDate(2021, 8, 1)},
		{"Maximum date boundary", NewDate(2127, 12, 31), maximumDate},
		{"Already at maximum", maximumDate, maximumDate},
		{"Minimum date", minimumDate, NewDate(2000, 1, 2)},
		{"Near minimum date", NewDate(2000, 1, 2), NewDate(2000, 1, 3)},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.NextDay(); got != tt.want {
				t.Errorf("NextDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousDay(t *testing.T) {
	var cases = []struct {
		name string
		this Date
		want Date
	}{
		{"Regular day end of month", NewDate(2021, 12, 31), NewDate(2021, 12, 30)},
		{"Regular day middle of month", NewDate(2021, 12, 2), NewDate(2021, 12, 1)},
		{"First day to previous month - Feb", NewDate(2021, 3, 1), NewDate(2021, 2, 28)},
		{"First day to leap day", NewDate(2024, 3, 1), NewDate(2024, 2, 29)}, // Leap year
		{"Regular consecutive days", NewDate(2021, 3, 2), NewDate(2021, 3, 1)},
		{"Year boundary", NewDate(2022, 1, 1), NewDate(2021, 12, 31)},       // Year boundary
		{"Already at minimum", minimumDate, minimumDate},                    // Already at minimum
		{"One day after minimum", NewDate(2000, 1, 2), NewDate(2000, 1, 1)}, // Near minimum
		{"Below minimum date", NewDate(1999, 1, 1), minimumDate},            // Below minimum
		{"First day of April", NewDate(2023, 4, 1), NewDate(2023, 3, 31)},   // Month with 31 days
		{"First day of month after 30-day month", NewDate(2023, 6, 1), NewDate(2023, 5, 31)},
		{"First day of March in leap year", NewDate(2020, 3, 1), NewDate(2020, 2, 29)},
		{"First day of March in regular year", NewDate(2023, 3, 1), NewDate(2023, 2, 28)},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.PreviousDay(); got != tt.want {
				t.Errorf("PreviousDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextWeek(t *testing.T) {
	cases := []struct {
		name       string
		this, want Date
	}{
		{"End of year boundary", NewDate(2021, 12, 31), NewDate(2022, 1, 7)},
		{"Start of month", NewDate(2021, 12, 1), NewDate(2021, 12, 8)},
		{"End of February non-leap year", NewDate(2021, 2, 28), NewDate(2021, 3, 7)},
		{"Start of February", NewDate(2021, 2, 1), NewDate(2021, 2, 8)},
		{"Leap day", NewDate(2024, 2, 29), NewDate(2024, 3, 7)},
		{"Near maximum date", NewDate(2127, 12, 25), NewDate(2128, 1, 1)},
		{"Last day of maximum year", NewDate(2127, 12, 31), maximumDate},
		{"Already at maximum date", maximumDate, maximumDate},
		{"Minimum date", minimumDate, NewDate(2000, 1, 8)},
		{"Near minimum date", NewDate(2000, 1, 2), NewDate(2000, 1, 9)},
		{"Month with 31 days", NewDate(2023, 5, 31), NewDate(2023, 6, 7)},
		{"Month with 30 days", NewDate(2023, 4, 30), NewDate(2023, 5, 7)},
		{"Mid-month regular case", NewDate(2023, 7, 15), NewDate(2023, 7, 22)},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.NextWeek(); got != tt.want {
				t.Errorf("NextWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousWeek(t *testing.T) {
	tests := []struct {
		name       string
		this, want Date
	}{
		{"Standard case - end of month", NewDate(2021, 12, 31), NewDate(2021, 12, 24)},
		{"Year boundary - early January", NewDate(2022, 1, 3), NewDate(2021, 12, 27)},
		{"Month boundary - start of month", NewDate(2021, 12, 1), NewDate(2021, 11, 24)},
		{"February to previous month", NewDate(2021, 3, 2), NewDate(2021, 2, 23)},
		{"Within same month", NewDate(2021, 3, 9), NewDate(2021, 3, 2)},
		{"Leap year February", NewDate(2024, 3, 5), NewDate(2024, 2, 27)},
		{"Near minimum date", NewDate(2000, 1, 8), NewDate(2000, 1, 1)},
		{"Exactly 7 days from minimum", NewDate(2000, 1, 8), NewDate(2000, 1, 1)},
		{"Less than 7 days from minimum", NewDate(2000, 1, 4), NewDate(2000, 1, 1)},
		{"At minimum edge", NewDate(2000, 1, 1), minimumDate},
		{"Already at minimum", minimumDate, minimumDate},
		{"Edge case - before minimum", NewDate(1999, 12, 31), minimumDate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.PreviousWeek(); got != tt.want {
				t.Errorf("PreviousWeek() = %v, want %v", got, tt.want)
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
		{"2000-01-01", minimumDate, false},
		{"2127-12-31", maximumDate, false},
		{"2024-02-29", NewDate(2024, 2, 29), false},
		{"invalid", Date(0), true},
		{"2021-13-01", NewDate(2022, 1, 1), false},   // Month 13 normalizes to next year
		{"2021-12-32", NewDate(2022, 1, 1), false},   // Day 32 normalizes to next month
		{"21-12-01", Date(0), true},                  // Wrong format
		{"2021/12/01", Date(0), true},                // Wrong separators
		{"", Date(0), true},                          // Empty string
		{"2021-1-1", Date(0), true},                  // Missing leading zeros
		{"2021-00-15", NewDate(2020, 12, 15), false}, // Month 0 normalizes
		{"2021-06-00", NewDate(2021, 5, 31), false},  // Day 0 normalizes
		{"2021-XX-15", Date(0), true},                // Invalid month characters
		{"2021-12-XX", Date(0), true},                // Invalid day characters
		{"XXXX-12-15", Date(0), true},                // Invalid year characters
		{"2127-12-32", maximumDate, false},           // Normalizes past maximum, clamped
		{"2000-01-00", minimumDate, false},           // Normalizes below minimum, clamped
	}
	for _, tt := range tests {
		t.Run(tt.formattedDate, func(t *testing.T) {
			got, err := Parse(tt.formattedDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
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
		{Date(0), "null"},
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

func TestDate_IsSet(t *testing.T) {
	tests := []struct {
		name     string
		thisDate Date
		want     bool
	}{
		{"Zero date", Date(0), false},
		{"Valid date", NewDate(2021, 12, 1), true},
		{"Minimum date", minimumDate, true},
		{"Maximum date", maximumDate, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.thisDate.IsSet(); got != tt.want {
				t.Errorf("IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MonthAfter(t *testing.T) {
	tests := []struct {
		name     string
		thisDate Date
		target   Date
		want     bool
	}{
		{"Same year, later month", NewDate(2021, 6, 15), NewDate(2021, 3, 10), true},
		{"Same year, earlier month", NewDate(2021, 3, 15), NewDate(2021, 6, 10), false},
		{"Same year, same month", NewDate(2021, 6, 15), NewDate(2021, 6, 10), false},
		{"Later year, earlier month", NewDate(2022, 1, 15), NewDate(2021, 12, 10), true},
		{"Earlier year, later month", NewDate(2021, 12, 15), NewDate(2022, 1, 10), false},
		{"Same date", NewDate(2021, 6, 15), NewDate(2021, 6, 15), false},
		{"One year later", NewDate(2022, 6, 15), NewDate(2021, 6, 15), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.thisDate.MonthAfter(tt.target); got != tt.want {
				t.Errorf("MonthAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MonthBefore(t *testing.T) {
	tests := []struct {
		name     string
		thisDate Date
		target   Date
		want     bool
	}{
		{"Same year, earlier month", NewDate(2021, 3, 15), NewDate(2021, 6, 10), true},
		{"Same year, later month", NewDate(2021, 6, 15), NewDate(2021, 3, 10), false},
		{"Same year, same month", NewDate(2021, 6, 15), NewDate(2021, 6, 10), false},
		{"Earlier year, later month", NewDate(2021, 12, 15), NewDate(2022, 1, 10), true},
		{"Later year, earlier month", NewDate(2022, 1, 15), NewDate(2021, 12, 10), false},
		{"Same date", NewDate(2021, 6, 15), NewDate(2021, 6, 15), false},
		{"One year earlier", NewDate(2021, 6, 15), NewDate(2022, 6, 15), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.thisDate.MonthBefore(tt.target); got != tt.want {
				t.Errorf("MonthBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Time(t *testing.T) {
	tests := []struct {
		name     string
		thisDate Date
	}{
		{"Regular date", NewDate(2021, 6, 15)},
		{"Minimum date", minimumDate},
		{"Maximum date", maximumDate},
		{"Leap day", NewDate(2024, 2, 29)},
		{"End of year", NewDate(2021, 12, 31)},
		{"Start of year", NewDate(2021, 1, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.thisDate.Time()
			if got.Year() != int(tt.thisDate.Year()) {
				t.Errorf("Time().Year() = %v, want %v", got.Year(), tt.thisDate.Year())
			}
			if byte(got.Month()) != tt.thisDate.Month() {
				t.Errorf("Time().Month() = %v, want %v", got.Month(), tt.thisDate.Month())
			}
			if got.Day() != int(tt.thisDate.Day()) {
				t.Errorf("Time().Day() = %v, want %v", got.Day(), tt.thisDate.Day())
			}
		})
	}
}

func TestNewDate_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		year        uint16
		month, day  byte
		wantYear    uint16
		wantMonth   byte
		wantDay     byte
		description string
	}{
		{"Zero month", 2021, 0, 15, 2020, 12, 15, "Month 0 should normalize"},
		{"Zero day", 2021, 6, 0, 2021, 5, 31, "Day 0 should normalize"},
		{"Month 13", 2021, 13, 15, 2022, 1, 15, "Month 13 should normalize to next year"},
		{"Day 32 in January", 2021, 1, 32, 2021, 2, 1, "Day 32 should normalize"},
		{"Day 30 in February non-leap", 2021, 2, 30, 2021, 3, 2, "Feb 30 should normalize"},
		{"Day 31 in April", 2021, 4, 31, 2021, 5, 1, "April 31 should normalize"},
		{"Month 13 at maximum year", 2127, 13, 1, 2127, 12, 31, "Normalization past maximum should clamp, not wrap"},
		{"Day 32 at maximum date", 2127, 12, 32, 2127, 12, 31, "Normalization past maximum should clamp, not wrap"},
		{"Day 0 at minimum date", 2000, 1, 0, 2000, 1, 1, "Normalization below minimum should clamp, not wrap"},
		{"Month 0 at minimum year", 2000, 0, 15, 2000, 1, 1, "Normalization below minimum should clamp, not wrap"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDate(tt.year, tt.month, tt.day)
			if got.Year() != tt.wantYear {
				t.Errorf("NewDate(%v, %v, %v).Year() = %v, want %v (%s)",
					tt.year, tt.month, tt.day, got.Year(), tt.wantYear, tt.description)
			}
			if got.Month() != tt.wantMonth {
				t.Errorf("NewDate(%v, %v, %v).Month() = %v, want %v (%s)",
					tt.year, tt.month, tt.day, got.Month(), tt.wantMonth, tt.description)
			}
			if got.Day() != tt.wantDay {
				t.Errorf("NewDate(%v, %v, %v).Day() = %v, want %v (%s)",
					tt.year, tt.month, tt.day, got.Day(), tt.wantDay, tt.description)
			}
		})
	}
}

func TestNow(t *testing.T) {
	now := Now()
	if !now.IsSet() {
		t.Errorf("Now() returned unset date")
	}
	currentTime := time.Now().UTC()
	if now.Year() != uint16(currentTime.Year()) {
		t.Errorf("Now().Year() = %v, want %v", now.Year(), currentTime.Year())
	}
}
