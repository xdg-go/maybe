package maybe

// I implements the Maybe monad for a int.  An I is considered 'valid' or
// 'invalid' depending on whether it contains a int or an error value.
type I struct {
	just int
	err  error
}

// NewI constructs an I from a given int or error. If e is not nil, returns
// ErrI(e), otherwise returns JustI(s)
func NewI(s int, e error) I {
	if e != nil {
		return ErrI(e)
	}
	return JustI(s)
}

// JustI constructs a valid I from a given int.
func JustI(s int) I {
	return I{just: s}
}

// ErrI constructs an invalid I from a given error.
func ErrI(e error) I {
	return I{err: e}
}

// IsErr returns true for an invalid I.
func (m I) IsErr() bool {
	return m.err != nil
}

// Bind applies a function that takes a int and returns an I.
func (m I) Bind(f func(s int) I) I {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// Split applies a function that takes a int and returns an AoI.
func (m I) Split(f func(s int) AoI) AoI {
	if m.err != nil {
		return ErrAoI(m.err)
	}

	return f(m.just)
}

// Unbox returns the underlying int value or error.
func (m I) Unbox() (int, error) {
	return m.just, m.err
}
