package wit

import (
	"sync"

	"github.com/manvalls/wit"
)

// Slice holds a slice of deltas internally
type Slice struct {
	aggregator *deltaAggregator
}

type deltaAggregator struct {
	sync.Mutex
	deltas []Delta
}

// NewSlice builds a new Slice
func NewSlice() Slice {
	return Slice{&deltaAggregator{sync.Mutex{}, []Delta{}}}
}

// Append appends deltas to the internal buffer
func (s Slice) Append(deltas ...Delta) {
	s.aggregator.Lock()
	defer s.aggregator.Unlock()
	s.aggregator.deltas = append(s.aggregator.deltas, deltas...)
}

// Delta flushes the internal buffer to the returned delta
func (s Slice) Delta() Delta {
	s.aggregator.Lock()
	defer s.aggregator.Unlock()

	delta := wit.List(s.aggregator.deltas...)
	s.aggregator.deltas = []Delta{}
	return delta
}
