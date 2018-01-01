package maybe

// S implements the Maybe monad for a string.  An S is considered 'valid' or
// 'invalid' depending on whether it contains a string or an error value.
type S struct {
	just string
	err  error
}

// NewS constructs an S from a given string or error. If e is not nil, returns
// ErrS(e), otherwise returns JustS(s)
func NewS(s string, e error) S {
	if e != nil {
		return ErrS(e)
	}
	return JustS(s)
}

// JustS constructs a valid S from a given string.
func JustS(s string) S {
	return S{just: s}
}

// ErrS constructs an invalid S from a given error.
func ErrS(e error) S {
	return S{err: e}
}

// IsErr returns true for an invalid S.
func (m S) IsErr() bool {
	return m.err != nil
}

// Bind applies a function that takes a string and returns an S.
func (m S) Bind(f func(s string) S) S {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// Split applies a function that takes a string and returns an AoS.
func (m S) Split(f func(s string) AoS) AoS {
	if m.err != nil {
		return ErrAoS(m.err)
	}

	return f(m.just)
}

// ToInt applies a function that takes a string and returns an I.
func (m S) ToInt(f func(s string) I) I {
	if m.err != nil {
		return ErrI(m.err)
	}

	return f(m.just)
}

// Unbox returns the underlying string value or error.
func (m S) Unbox() (string, error) {
	return m.just, m.err
}
