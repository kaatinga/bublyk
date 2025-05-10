package bublyk

import "testing"

func BenchmarkNewString(b *testing.B) {
	b.ReportAllocs()
	date := Now()
	for i := 0; i < b.N; i++ {
		_ = date.String()
	}
}

// func BenchmarkOldString(b *testing.B) {
//	b.ReportAllocs()
//	date := Now()
//	for i := 0; i < b.N; i++ {
//		date.oldString()
//	}
// }

func BenchmarkWithAssets(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		maximumDate.Format(postgreSQLFormat)
	}
}

func BenchmarkUsingTime(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		maximumDate.Format("2006_01-02")
	}
}

func BenchmarkShiftByTime(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		NewDate(2021, 12, 31).NextDay()
	}
}

func BenchmarkShiftByBitwise(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		NewDate(2021, 12, 25).NextDay()
	}
}
