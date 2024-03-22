package service

import (
	"context"
	"fmt"
	"github.com/eldius/1brc-go/internal/decoder"
	"github.com/eldius/1brc-go/internal/reader"
	"golang.org/x/sync/errgroup"
	"sync"
)

func Consume(queueSize, workersCount int, fileName string) {
	in := make(chan [2]string, queueSize)
	go func(in chan [2]string) {
		if err := reader.Read(fileName, in); err != nil {
			err = fmt.Errorf("setting up readers: %w", err)
			panic(err)
		}
	}(in)

	var wg sync.WaitGroup
	wg.Add(workersCount)
	p := decoder.NewProcessorWg()
	for _ = range workersCount {
		go p.Process(&wg, in)
	}
	wg.Wait()

	p.Print()
}

func ConsumeAlt(queueSize, workersCount int, fileName string) error {
	in := make(chan [2]string, queueSize)
	go func(in chan [2]string) {
		if err := reader.Read(fileName, in); err != nil {
			err = fmt.Errorf("setting up readers: %w", err)
			panic(err)
		}
	}(in)

	eg, _ := errgroup.WithContext(context.Background())
	p := decoder.NewProcessorEg()
	for _ = range workersCount {
		eg.Go(func() error {
			return p.Process(in)
		})
	}

	if err := eg.Wait(); err != nil {
		err = fmt.Errorf("waiting to tasks to finish: %w", err)
		return err
	}
	p.Print()
	return nil
}
