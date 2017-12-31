package maybe

// S implements the Maybe monad for a string.
type S struct {
	just string
	err  error
}

// NewS constructs an S from a given slice of strings or error. If e is
// not nil, returns ErrS(e), otherwise returns JustS(s)
func NewS(s string, e error) S {
	if e != nil {
		return ErrS(e)
	}
	return JustS(s)
}

// JustS constructs a "Just" S from a given slice of strings.
func JustS(s string) S {
	return S{just: s}
}

// ErrS constructs a "Nothing" S from a given error.
func ErrS(e error) S {
	return S{err: e}
}

// IsErr returns true for a "Nothing" S with an error
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

// Unbox returns the underlying string value or error.
func (m S) Unbox() (string, error) {
	return m.just, m.err
}
