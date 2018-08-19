package wit

import (
	"context"
	"sync"
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

// RunAppend runs and appends to the internal buffer the given function
func (s Slice) RunAppend(parentCtx context.Context, callback func(context.Context) Delta) {
	s.Append(Run(parentCtx, callback))
}

// RunAppendSync runs and appends to the internal buffer the given function synchronously
func (s Slice) RunAppendSync(callback func() Delta) {
	s.Append(RunSync(callback))
}

// Delta flushes the internal buffer to the returned delta
func (s Slice) Delta() Delta {
	s.aggregator.Lock()
	defer s.aggregator.Unlock()

	delta := List(s.aggregator.deltas...)
	s.aggregator.deltas = []Delta{}
	return delta
}
