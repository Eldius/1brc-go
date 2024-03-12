package main

import (
	"flag"
	"fmt"
	"github.com/eldius/1brc-go/internal/parser"
	"github.com/eldius/1brc-go/internal/process"
	"log/slog"
	"os"
	"sync"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})))
}

func main() {
	fileName := flag.String("file", "./internal/parser/sample_data/measurements_50.txt", "File to be parsed")
	workersCount := flag.Int("workers-count", 5, "File to be parsed")
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
		if err := parser.Read(*fileName, in); err != nil {
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
