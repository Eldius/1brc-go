package process

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"
)

type Record struct { //32b
	min   float64 //8b
	max   float64 //8b
	count int64   //8b
	mean  float64 //8b
}

func (r *Record) Process(d string) {
	m, _ := strconv.ParseFloat(d, 64)
	if r.min > m {
		r.min = m
	}
	if r.max < m {
		r.max = m
	}
	r.mean = ((r.mean * float64(r.count)) + m) / float64(r.count+1)
	r.count++
}

type ProcessedData struct {
	d map[string]*Record
	m sync.Mutex
}

func (p *ProcessedData) Add(d [2]string) {
	p.m.Lock()
	defer p.m.Unlock()
	if data, ok := p.d[d[0]]; ok {
		data.Process(d[1])
	} else {
		r := newRecord()
		r.Process(d[1])
		p.d[d[0]] = r
	}
}

func (p *Processor) Print() {
	for k, v := range p.d.d {
		fmt.Printf("%s: %01.1f/%01.1f/%01.1f\n", k, v.min, v.mean, v.max)
	}
}

func newProcessData() *ProcessedData {
	return &ProcessedData{d: make(map[string]*Record)}
}

func newRecord() *Record {
	r := Record{
		min:   1.7e+308,
		max:   2.2e-308,
		count: 0,
		mean:  0,
	}

	return &r
}

type Processor struct {
	d *ProcessedData
}

func (p *Processor) Process(wg *sync.WaitGroup, ch chan [2]string) {
	defer wg.Done()
	for d := range ch {
		slog.With(slog.String("line", fmt.Sprintf("%v", d))).Info("processing")
		p.d.Add(d)
	}
}

func (p *Processor) Data() *ProcessedData {
	return p.d
}

func NewProcessor() *Processor {
	return &Processor{d: newProcessData()}
}
