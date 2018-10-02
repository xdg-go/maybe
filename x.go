package maybe

import "fmt"

// X implements the Maybe monad for an empty interface.  An X is considered
// 'valid' or 'invalid' depending on whether it contains a non-nil interface
// or an error value.
type X struct {
	just interface{}
	err  error
}

// NewX constructs an X from a given empty interface or error. If e is not nil, returns
// ErrX(e), otherwise returns JustX(s)
func NewX(x interface{}, e error) X {
	if e != nil {
		return ErrX(e)
	}
	return JustX(x)
}

// JustX constructs a valid X from a given empty interface.
func JustX(x interface{}) X {
	return X{just: x}
}

// ErrX constructs an invalid X from a given error.
func ErrX(e error) X {
	return X{err: e}
}

// IsErr returns true for an invalid X.
func (m X) IsErr() bool {
	return m.just == nil || m.err != nil
}

// Bind applies a function that takes an interface and returns an X.
func (m X) Bind(f func(x interface{}) X) X {
	if m.IsErr() {
		return m
	}

	return f(m.just)
}

// Split applies a function that takes an interface and returns an AoX.
func (m X) Split(f func(x interface{}) AoX) AoX {
	if m.IsErr() {
		return ErrAoX(m.err)
	}

	return f(m.just)
}

// String returns a string representation, mostly useful for debugging.
func (m X) String() string {
	if m.IsErr() {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// Unbox returns the underlying empty interface value or error.
func (m X) Unbox() (interface{}, error) {
	return m.just, m.err
}
