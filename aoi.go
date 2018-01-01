package maybe

// AoI implements the Maybe monad for a slice of ints.
type AoI struct {
	just []int
	err  error
}

// NewAoI constructs an AoI from a given slice of ints or error. If e is
// not nil, returns ErrAOS(e), otherwise returns JustAOS(s)
func NewAoI(s []int, e error) AoI {
	if e != nil {
		return ErrAoI(e)
	}
	return JustAoI(s)
}

// JustAoI constructs a "Just" AoI from a given slice of ints.
func JustAoI(s []int) AoI {
	return AoI{just: s}
}

// ErrAoI constructs a "Nothing" AoI from a given error.
func ErrAoI(e error) AoI {
	return AoI{err: e}
}

// IsErr returns true for a "Nothing" AoI with an error
func (m AoI) IsErr() bool {
	return m.err != nil
}

// Bind applies a function that takes a slice of ints and returns an AoI.
func (m AoI) Bind(f func(s []int) AoI) AoI {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a slice of ints and returns an I.
func (m AoI) Join(f func(s []int) I) I {
	if m.err != nil {
		return ErrI(m.err)
	}

	return f(m.just)
}

// Map applies a function to each element of a valid AoI and returns a new
// AoI.  If the AoI is invalid or if any function returns an invalid I, Map
// returns a "Nothing" AoI.
func (m AoI) Map(f func(s int) I) AoI {
	if m.err != nil {
		return m
	}

	new := make([]int, len(m.just))
	for i, v := range m.just {
		x, err := f(v).Unbox()
		if err != nil {
			return ErrAoI(err)
		}
		new[i] = x
	}

	return JustAoI(new)
}

// Unbox returns the underlying slice of ints value or error.
func (m AoI) Unbox() ([]int, error) {
	return m.just, m.err
}
