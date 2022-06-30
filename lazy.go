package booklit

import "sync"

type Lazy struct {
	Block bool

	fn     func() (Content, error)
	result Content
	err    error
	once   *sync.Once
}

func LazyFlow(fn func() (Content, error)) *Lazy {
	return &Lazy{
		Block: false,
		fn:    fn,
		once:  &sync.Once{},
	}
}

func LazyBlock(fn func() (Content, error)) *Lazy {
	return &Lazy{
		Block: true,
		fn:    fn,
		once:  &sync.Once{},
	}
}

// IsFlow returns the value of Block and never delegates to the deferred
// content, i.e. this property must always be known ahead of time.
func (lazy *Lazy) IsFlow() bool {
	return !lazy.Block
}

// String returns the string value.
func (lazy *Lazy) String() string {
	if lazy.Block {
		return "<lazy block>"
	} else {
		return "<lazy flow>"
	}
}

// Visit calls VisitString.
func (lazy *Lazy) Visit(visitor Visitor) error {
	return visitor.VisitLazy(lazy)
}

// Force generates the deferred content if needed and returns it.
func (lazy *Lazy) Force() (Content, error) {
	lazy.once.Do(func() {
		lazy.result, lazy.err = lazy.fn()
	})

	return lazy.result, lazy.err
}
