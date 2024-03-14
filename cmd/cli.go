package main

import (
	"flag"
	"fmt"
	"github.com/eldius/1brc-go/internal/process"
	"github.com/eldius/1brc-go/internal/reader"
	"log/slog"
	"os"
	"slices"
	"sync"
)

func init() {
	hostname, _ := os.Hostname()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if slices.Contains([]string{"host.name", "service.name", "level", "message"}, a.Key) {
				return a
			}
			if a.Key == "msg" {
				a.Key = "message"
				return a
			}
			a.Key = "custom.1brc." + a.Key
			return a
		},
	})).With(
		slog.String("hostname", hostname),
		slog.String("service.name", "1brc"),
	))
}

func main() {
	fileName := flag.String("file", "measurements.txt", "File to be parsed")
	workersCount := flag.Int("workers-count", 5, "Record processors count")
	queueSize := flag.Int("queue-size", 5, "Process queue size")

	flag.Parse()

	log := slog.With(
		slog.String("file", *fileName),
		slog.Int("workers-count", *workersCount),
		slog.Int("queue-size", *workersCount),
	)

	log.Info("starting...")

	in := make(chan [2]string, *queueSize)
	go func(in chan [2]string) {
		if err := reader.Read(*fileName, in); err != nil {
			err = fmt.Errorf("setting up readers: %w", err)
			panic(err)
		}
	}(in)

	var wg sync.WaitGroup
	wg.Add(*workersCount)
	p := process.NewProcessor()
	for _ = range *workersCount {
		go p.Process(&wg, in)
	}
	wg.Wait()

	p.Print()
}
