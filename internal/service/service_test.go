package service

import "testing"

func BenchmarkProcessFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Consume(10, 5, "./service/test_data/measurements_5.txt")
		Consume(10, 5, "./service/test_data/measurements_50.txt")
	}
}
