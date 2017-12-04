package maybe

// S implements the Maybe monad for a string.
type S struct {
	just string
	err  error
}

// NewS constructs a "Just" S from a given string.
func NewS(s string) S {
	return S{just: s}
}

// ErrS constructs a "Nothing" S from a given error.
func ErrS(e error) S {
	return S{err: e}
}

// Bind applies a function that takes a string and returns an S.
func (m S) Bind(f func(s string) S) S {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// BindAoS applies a function that takes a string and returns an AoS.
func (m S) BindAoS(f func(s string) AoS) AoS {
	if m.err != nil {
		return ErrAoS(m.err)
	}

	return f(m.just)
}

// Unbox returns the underlying string value or error.
func (m S) Unbox() (string, error) {
	return m.just, m.err
}
