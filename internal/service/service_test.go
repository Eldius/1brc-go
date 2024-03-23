package service

import (
	"path/filepath"
	"testing"
)

func BenchmarkConsumeWG(b *testing.B) {
	runFileWg(b, "measurements_5.txt")
	runFileWg(b, "measurements_50.txt")
	runFileWg(b, "measurements_1k.txt")
}

func BenchmarkConsumeEG(b *testing.B) {
	runFileEG(b, "measurements_5.txt")
	runFileEG(b, "measurements_50.txt")
	runFileEG(b, "measurements_1k.txt")
}

func runFileWg(b *testing.B, file string) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConsumeWaitGroup(10, 5, filepath.Join("../service/test_data", file))
	}
}

func runFileEG(b *testing.B, file string) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := ConsumeErrorGroup(10, 5, filepath.Join("../service/test_data", file)); err != nil {
			b.FailNow()
		}
	}
}
