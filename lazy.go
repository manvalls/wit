package wit

import "sync"

type lazyAction struct {
	f     func() Action
	mux   sync.Mutex
	delta *Delta
}

func (l *lazyAction) Delta() Delta {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.delta != nil {
		return *l.delta
	}

	delta := l.f().Delta()
	l.delta = &delta
	return delta
}

// Lazy builds an action which will be evaluated on demand
func Lazy(f func() Action) Action {
	return &lazyAction{f: f}
}

// Async spanws the provided function in a new goroutine and
// returns an action which resolves to the returned one
func Async(f func() Action) Action {
	l := Lazy(f)
	go l.Delta()
	return l
}

// Memo returns an action which will be memoized
func Memo(a Action) Action {
	return Lazy(func() Action {
		return a
	})
}
