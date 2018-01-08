package maybe

import "fmt"

// AoAoI implements the Maybe monad for a 2D slice of ints.  An AoAoI is considered
// 'valid' or 'invalid' depending on whether it contains a 2D slice of ints or an
// error value.
type AoAoI struct {
	just [][]int
	err  error
}

// NewAoAoI constructs an AoAoI from a given 2D slice of ints or error. If e is not
// nil, returns ErrAoAoI(e), otherwise returns JustAoAoI(s).
func NewAoAoI(s [][]int, e error) AoAoI {
	if e != nil {
		return ErrAoAoI(e)
	}
	return JustAoAoI(s)
}

// JustAoAoI constructs a valid AoAoI from a given 2D slice of ints.
func JustAoAoI(s [][]int) AoAoI {
	return AoAoI{just: s}
}

// ErrAoAoI constructs an invalid AoAoI from a given error.
func ErrAoAoI(e error) AoAoI {
	return AoAoI{err: e}
}

// IsErr returns true for an invalid AoAoI.
func (m AoAoI) IsErr() bool {
	return m.err != nil
}

// Bind applies a function that takes a 2D slice of ints and returns an AoAoI.
func (m AoAoI) Bind(f func(s [][]int) AoAoI) AoAoI {
	if m.err != nil {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a 2D slice of ints and returns an AoI.
func (m AoAoI) Join(f func(s [][]int) AoI) AoI {
	if m.err != nil {
		return ErrAoI(m.err)
	}

	return f(m.just)
}

// Map applies a function to each element of a valid AoAoI and returns a new
// AoAoI.  If the AoAoI is invalid or if any function returns an invalid I,
// Map returns an invalid AoAoI.
func (m AoAoI) Map(f func(s []int) AoI) AoAoI {
	if m.err != nil {
		return m
	}

	new := make([][]int, len(m.just))
	for i, v := range m.just {
		x, err := f(v).Unbox()
		if err != nil {
			return ErrAoAoI(err)
		}
		new[i] = x
	}

	return JustAoAoI(new)
}

// String returns a string representation, mostly useful for debugging.
func (m AoAoI) String() string {
	if m.err != nil {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// ToStr applies a function that takes a string and returns an S.  If the
// AoAoI is invalid or if any function returns an invalid S, ToStr returns an
// invalid AoS.  Note: unlike Map, this is a deep conversion of individual
// elements of the 2D slice of ints.
func (m AoAoI) ToStr(f func(x int) S) AoAoS {
	if m.err != nil {
		return ErrAoAoS(m.err)
	}

	new := make([][]string, len(m.just))
	for i, xs := range m.just {
		new[i] = make([]string, len(xs))
		for j, v := range xs {
			num, err := f(v).Unbox()
			if err != nil {
				return ErrAoAoS(err)
			}
			new[i][j] = num
		}
	}

	return JustAoAoS(new)
}

// Unbox returns the underlying 2D slice of ints or error.
func (m AoAoI) Unbox() ([][]int, error) {
	return m.just, m.err
}
