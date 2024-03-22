package decoder

import (
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
	"unsafe"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestProcess(t *testing.T) {
	t.Run("single location with two values", func(t *testing.T) {
		p := NewProcessorWg()
		var wg sync.WaitGroup
		wg.Add(5)
		in := make(chan [2]string)
		go p.Process(&wg, in)

		in <- [2]string{"test0", "10.0"}
		in <- [2]string{"test0", "5.0"}
		in <- [2]string{"test0", "10.0"}
		in <- [2]string{"test0", "1.0"}
		in <- [2]string{"test0", "10.7"}

		r := p.Data().d["test0"]

		if !almostEqual(r.min, 1.0) {
			t.Logf("min is not 1.0, was %0.1f", r.min)
			t.Fail()
		}

		if !almostEqual(r.max, 10.7) {
			t.Logf("min is not 10.7, was %0.1f", r.max)
			t.Fail()
		}

		if !almostEqual(r.mean, 7.34) {
			t.Logf("min is not 7.34, was %0.1f", r.mean)
			t.Fail()
		}
	})
}

func TestName(t *testing.T) {
	r := newRecord()
	t.Logf("record struct size: %d", unsafe.Sizeof(*r))
	d := newDecodedData()

	t.Logf("[before] processed data struct size: %d", unsafe.Sizeof(*d))

	assert.Nil(t, d.Add([2]string{"Rio", "1.5"}))
	assert.Nil(t, d.Add([2]string{"Sampa", "1.5"}))
	assert.Nil(t, d.Add([2]string{"Rio", "2.5"}))
	assert.Nil(t, d.Add([2]string{"Sampa", "3.5"}))
	assert.Nil(t, d.Add([2]string{"Minas", "3.5"}))
	assert.Nil(t, d.Add([2]string{"Minas", "3.5"}))
	assert.Nil(t, d.Add([2]string{"Sampa", "3.5"}))
	assert.Nil(t, d.Add([2]string{"Sampa", "3.5"}))

	t.Logf("[after] processed data struct size: %d", unsafe.Sizeof(*d))
}
