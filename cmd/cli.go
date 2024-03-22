package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/eldius/1brc-go/internal/process"
	"github.com/eldius/1brc-go/internal/reader"
	"golang.org/x/exp/trace"
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
	traceEnabled := flag.Bool("trace", false, "Enable trace recording")

	flag.Parse()

	log := slog.With(
		slog.String("file", *fileName),
		slog.Int("workers-count", *workersCount),
		slog.Int("queue-size", *workersCount),
	)

	deferredFunc := func() func() {
		if *traceEnabled {
			fr := trace.NewFlightRecorder()
			_ = fr.Start()

			return func() {
				var b bytes.Buffer
				_, err := fr.WriteTo(&b)
				if err != nil {
					err = fmt.Errorf("parsing trace data: %w", err)
					log.With("error", err).Error("parsing trace data")
					return
				}
				// Write it to a file.
				if err := os.WriteFile("trace.out", b.Bytes(), 0o755); err != nil {
					err = fmt.Errorf("writing trace data to file: %w", err)
					log.With("error", err).Error("writing trace data to file")
					return
				}
			}
		}
		return func() {}
	}()
	defer deferredFunc()

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
