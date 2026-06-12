[![Tests and linter](https://github.com/kaatinga/bublyk/actions/workflows/golang_ci.yml/badge.svg)](https://github.com/kaatinga/bublyk/actions/workflows/golang_ci.yml)
[![GitHub release](https://img.shields.io/github/release/kaatinga/bublyk.svg)](https://github.com/kaatinga/bublyk/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/bublyk/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/bublyk/branch/main/graph/badge.svg?token=Q34SE0KN9E)](https://codecov.io/gh/kaatinga/bublyk)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/bublyk/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# bublyk

A lightweight, memory-efficient Go package for working with dates without time components. The `bublyk` package introduces the `Date` type, specifically designed for applications that only need calendar dates in UTC, without the overhead of time details.

## Features

### Memory Efficiency
- **Only 2 bytes per date** - Uses `uint16` internally with bit packing
- **87.5% less memory** than `time.Time` (16 bytes)
- Perfect for large datasets, databases, and memory-constrained environments

### Developer Experience
- **Direct comparisons** - Use standard operators (`>`, `<`, `==`, `>=`, `<=`)
- **No boilerplate** - Clean, intuitive API
- **Automatic normalization** - Invalid dates are automatically corrected
- **100% test coverage** - Thoroughly tested and reliable

### Date Range
- **Supported range**: January 1, 2000 to December 31, 2127
- Dates outside this range are clamped to boundaries
- Optimized for modern applications

## Installation

```bash
go get github.com/kaatinga/bublyk
```

## Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/kaatinga/bublyk"
)

func main() {
	// Get current date
	today := bublyk.Now()
	fmt.Println("Today:", today) // Output: 2025-11-16

	// Create specific dates
	date1 := bublyk.NewDate(2024, 12, 25)
	date2 := bublyk.NewDate(2025, 1, 1)

	// Direct comparisons
	if date2 > date1 {
		fmt.Println("New Year comes after Christmas")
	}

	// Convert from time.Time
	t := time.Now().AddDate(0, 0, 5)
	futureDate := bublyk.NewDateFromTime(&t)
	fmt.Println("5 days from now:", futureDate)
}
```

## API Reference

### Creating Dates

#### `Now() Date`
Returns the current date in UTC.

```go
today := bublyk.Now()
```

#### `NewDate(year uint16, month, day byte) Date`
Creates a new date. Automatically normalizes invalid dates.

```go
date := bublyk.NewDate(2024, 2, 29) // Leap day
normalized := bublyk.NewDate(2024, 13, 1) // Becomes 2025-01-01
```

#### `NewDateFromTime(t *time.Time) Date`
Converts a `time.Time` to a `Date`. Returns zero date if `t` is nil.

```go
t := time.Now()
date := bublyk.NewDateFromTime(&t)
```

#### `CurrentMonth() Date`
Returns the first day of the current month.

```go
monthStart := bublyk.CurrentMonth() // e.g., 2025-11-01
```

#### `Parse(formattedDate string) (Date, error)`
Parses a date string in `YYYY-MM-DD` format.

```go
date, err := bublyk.Parse("2024-12-25")
if err != nil {
	log.Fatal(err)
}
```

### Date Methods

#### Accessors

```go
date := bublyk.NewDate(2024, 12, 25)

year := date.Year()   // uint16: 2024
month := date.Month() // byte: 12
day := date.Day()     // byte: 25
```

#### Checking and Comparing

```go
// Check if date is set (non-zero)
if date.IsSet() {
	fmt.Println("Date is valid")
}

// Check if date is in the future
if date.IsFuture() {
	fmt.Println("This date hasn't happened yet")
}

// Month-level comparisons
date1 := bublyk.NewDate(2024, 6, 15)
date2 := bublyk.NewDate(2024, 3, 20)

if date1.MonthAfter(date2) {
	fmt.Println("date1 is at least one month after date2")
}

if date2.MonthBefore(date1) {
	fmt.Println("date2 is at least one month before date1")
}
```

#### Date Navigation

```go
date := bublyk.NewDate(2024, 12, 25)

// Move by days
tomorrow := date.NextDay()
yesterday := date.PreviousDay()

// Move by weeks
nextWeek := date.NextWeek()
lastWeek := date.PreviousWeek()
```

#### Formatting

```go
date := bublyk.NewDate(2024, 12, 25)

// Default format (YYYY-MM-DD)
str := date.String() // "2024-12-25"

// Custom format using time.Time layouts
formatted := date.Format(time.RFC822) // "25 Dec 24 00:00 UTC"

// European format (DD.MM.YYYY)
european := date.DMYWithDots() // "25.12.2024"
```

#### Conversion

```go
date := bublyk.NewDate(2024, 12, 25)

// Convert to time.Time
t := date.Time() // *time.Time in UTC
```

## Use Cases

### Database Storage
Store dates efficiently in databases:

```go
type Event struct {
	ID        int
	Name      string
	EventDate bublyk.Date // Only 2 bytes!
}
```

### Date Ranges and Filtering

```go
startDate := bublyk.NewDate(2024, 1, 1)
endDate := bublyk.NewDate(2024, 12, 31)

for date := startDate; date <= endDate; date = date.NextDay() {
	// Process each day of 2024
}
```

### Business Logic

```go
// Check if subscription is expired
subscription := getSubscription()
if subscription.ExpiryDate < bublyk.Now() {
	sendRenewalNotice()
}

// Calculate billing periods (month 13 normalizes to January of the next year)
billingDate := bublyk.NewDate(2024, 1, 1)
for i := 0; i < 12; i++ {
	processBilling(billingDate)
	billingDate = bublyk.NewDate(billingDate.Year(), billingDate.Month()+1, 1)
}
```

## Technical Details

### Memory Layout

The `Date` type uses bit packing to store year, month, and day in a single `uint16`:

```
Bits: [15-9: Year offset] [8-5: Month] [4-0: Day]
      7 bits (0-127)      4 bits (0-15) 5 bits (0-31)
```

- **Year**: Stored as offset from 2000 (0-127 = years 2000-2127)
- **Month**: 4 bits (1-12)
- **Day**: 5 bits (1-31)

### Date Normalization

Invalid dates are automatically normalized using Go's `time` package:

```go
bublyk.NewDate(2024, 2, 30)  // → 2024-03-01
bublyk.NewDate(2024, 13, 1)  // → 2025-01-01
bublyk.NewDate(2024, 4, 31)  // → 2024-05-01
bublyk.NewDate(1999, 1, 1)   // → 2000-01-01 (minimum)
bublyk.NewDate(2200, 1, 1)   // → 2127-12-31 (maximum)
```

### Performance

- **Comparisons**: Direct integer comparison (extremely fast)
- **Memory**: 2 bytes vs 16 bytes for `time.Time`
- **Parsing**: Optimized with `faststrconv` package
- **Navigation**: Bit manipulation for same-month operations

## Limitations

- **Date range**: 2000-01-01 to 2127-12-31 only
- **No time component**: Only stores calendar dates
- **UTC only**: No timezone support (by design)
- **No locale**: Formatting is limited to standard layouts

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

MIT License - see [LICENSE](LICENSE) file for details.
