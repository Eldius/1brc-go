package decoder

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

func (r *Record) Decode(d string) error {
	m, err := strconv.ParseFloat(d, 64)
	if err != nil {
		err = fmt.Errorf("parsing data ('%s'): %w", d, err)
		return err
	}
	if r.min > m {
		r.min = m
	}
	if r.max < m {
		r.max = m
	}
	r.mean = ((r.mean * float64(r.count)) + m) / float64(r.count+1)
	r.count++

	return nil
}

type DecodedData struct {
	d map[string]*Record
	m sync.Mutex
}

func (p *DecodedData) Add(d [2]string) error {
	p.m.Lock()
	defer p.m.Unlock()
	if data, ok := p.d[d[0]]; ok {
		if err := data.Decode(d[1]); err != nil {
			err = fmt.Errorf("update existing location ('%s'): %w", d[0], err)
			return err
		}
	} else {
		r := newRecord()
		if err := r.Decode(d[1]); err != nil {
			err = fmt.Errorf("update new location ('%s'): %w", d[0], err)
			return err
		}
		p.d[d[0]] = r
	}

	return nil
}

func newDecodedData() *DecodedData {
	return &DecodedData{d: make(map[string]*Record)}
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

type ProcessorWg struct {
	d *DecodedData
}

func (p *ProcessorWg) Print() {
	for k, v := range p.d.d {
		fmt.Printf("%s: %01.1f/%01.1f/%01.1f\n", k, v.min, v.mean, v.max)
	}
}

func (p *ProcessorWg) Process(wg *sync.WaitGroup, ch chan [2]string) {
	defer wg.Done()
	for d := range ch {
		//slog.With(slog.String("line", fmt.Sprintf("%v", d))).Info("processing")
		_ = p.d.Add(d)
	}
}

func (p *ProcessorWg) Data() *DecodedData {
	return p.d
}

type ProcessorEg struct {
	d *DecodedData
}

func (p *ProcessorEg) Process(ch chan [2]string) error {
	for d := range ch {
		slog.With(slog.String("line", fmt.Sprintf("%v", d))).Info("processing")
		if err := p.d.Add(d); err != nil {
			err = fmt.Errorf("alternatively processing data ('%s'): %w", d, err)
			return err
		}
	}
	return nil
}

func (p *ProcessorEg) Print() {
	for k, v := range p.d.d {
		fmt.Printf("%s: %01.1f/%01.1f/%01.1f\n", k, v.min, v.mean, v.max)
	}
}

func (p *ProcessorEg) Data() *DecodedData {
	return p.d
}

func NewProcessorWg() *ProcessorWg {
	return &ProcessorWg{d: newDecodedData()}
}

func NewProcessorEg() *ProcessorEg {
	return &ProcessorEg{d: newDecodedData()}
}
