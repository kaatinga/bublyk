[![Tests](https://github.com/kaatinga/bublyk/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kaatinga/bublyk/actions/workflows/test.yml)
[![GitHub release](https://img.shields.io/github/release/kaatinga/bublyk.svg)](https://github.com/kaatinga/bublyk/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/bublyk/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/bublyk/branch/main/graph/badge.svg?token=Q34SE0KN9E)](https://codecov.io/gh/kaatinga/bublyk)
[![lint workflow](https://github.com/kaatinga/bublyk/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/bublyk/actions?query=workflow%3Alinter)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/bublyk/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# bublyk

The package introduces the Date type, specifically designed for instances where only the date in UTC location is required, without the need for time details. In comparison to the time.Time type, Date offers several advantages:
- It consumes significantly less memory.
- It eliminates the need for boilerplate code when working with dates.
- It allows for straightforward comparisons using operators such as >, <, and others.

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/kaatinga/bublyk"
)

func main() {
	// Create a new Date instance
	date := bublyk.Now()

	// Print the current date
	fmt.Println(date)

	// Create a new Date instance from a time.Time instance
	var t = time.Now().AddDate(0, 0, 5)
	date2 := bublyk.NewDateFromTime(&t)

	// Print the current date
	fmt.Println(date2)

	// Compare two Date instances
	fmt.Println(date > date2)
}

```

Will be happy to everyone who want to participate in the work on the `bublyk` package.
