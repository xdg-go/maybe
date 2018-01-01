package maybe

// AoS implements the Maybe monad for a slice of strings.
type AoS struct {
	just []string
	err  error
}

// NewAoS constructs an AoS from a given slice of strings or error. If e is
// not nil, returns ErrAOS(e), otherwise returns JustAOS(s)
func NewAoS(s []string, e error) AoS {
	if e != nil {
		return ErrAoS(e)
	}
	return JustAoS(s)
}

// JustAoS constructs a "Just" AoS from a given slice of strings.
func JustAoS(s []string) AoS {
	return AoS{just: s}
}

// ErrAoS constructs a "Nothing" AoS from a given error.
func ErrAoS(e error) AoS {
	return AoS{err: e}
}

// IsErr returns true for a "Nothing" AoS with an error
func (m AoS) IsErr() bool {
	return m.err != nil
}

// Bind applies a function that takes a slice of strings and returns an AoS.
func (m AoS) Bind(f func(s []string) AoS) AoS {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a slice of strings and returns an S.
func (m AoS) Join(f func(s []string) S) S {
	if m.err != nil {
		return ErrS(m.err)
	}

	return f(m.just)
}

// Map applies a function to each element of a valid AoS and returns a new
// AoS.  If the AoS is invalid or if any function returns an invalid S, Map
// returns a "Nothing" AoS.
func (m AoS) Map(f func(s string) S) AoS {
	if m.err != nil {
		return m
	}

	new := make([]string, len(m.just))
	for i, v := range m.just {
		str, err := f(v).Unbox()
		if err != nil {
			return ErrAoS(err)
		}
		new[i] = str
	}

	return JustAoS(new)
}

// Unbox returns the underlying slice of strings value or error.
func (m AoS) Unbox() ([]string, error) {
	return m.just, m.err
}
