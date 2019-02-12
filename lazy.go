package wit

import "sync"

type lazyCommand struct {
	f     func() Command
	mux   sync.Mutex
	delta *Delta
}

func (l *lazyCommand) Delta() Delta {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.delta != nil {
		return *l.delta
	}

	delta := l.f().Delta()
	l.delta = &delta
	return delta
}

// Lazy builds an command which will be evaluated on demand
func Lazy(f func() Command) Command {
	return &lazyCommand{f: f}
}

// Async spanws the provided function in a new goroutine and
// returns an command which resolves to the returned one
func Async(f func() Command) Command {
	l := Lazy(f)
	go l.Delta()
	return l
}

// Memo returns an command which will be memoized
func Memo(a Command) Command {
	return Lazy(func() Command {
		return a
	})
}
