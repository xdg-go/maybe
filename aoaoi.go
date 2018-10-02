package maybe

import (
	"errors"
	"fmt"
)

// AoAoI implements the Maybe monad for a 2-D slice of ints.  An AoAoI is
// considered 'valid' or 'invalid' depending on whether it contains a 2-D
// slice of ints or an error value.  A zero-value AoAoI is invalid and Unbox()
// will return an error to that effect.
type AoAoI struct {
	just [][]int
	err  error
}

// NewAoAoI constructs an AoAoI from a given 2-D slice of ints or error. If e is not
// nil, returns ErrAoAoI(e), otherwise returns JustAoAoI(s).
func NewAoAoI(s [][]int, e error) AoAoI {
	if e != nil {
		return ErrAoAoI(e)
	}
	return JustAoAoI(s)
}

// JustAoAoI constructs a valid AoAoI from a given 2-D slice of ints.
func JustAoAoI(s [][]int) AoAoI {
	return AoAoI{just: s}
}

// ErrAoAoI constructs an invalid AoAoI from a given error.
func ErrAoAoI(e error) AoAoI {
	return AoAoI{err: e}
}

// IsErr returns true for an invalid AoAoI.
func (m AoAoI) IsErr() bool {
	return m.just == nil || m.err != nil
}

// Bind applies a function that takes a 2-D slice of ints and returns an AoAoI.
func (m AoAoI) Bind(f func(s [][]int) AoAoI) AoAoI {
	if m.IsErr() {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a 2-D slice of ints and returns an AoI.
func (m AoAoI) Join(f func(s []int) I) AoI {
	if m.IsErr() {
		return ErrAoI(m.err)
	}

	new := make([]int, len(m.just))
	for i, v := range m.just {
		s, err := f(v).Unbox()
		if err != nil {
			return ErrAoI(err)
		}
		new[i] = s
	}

	return JustAoI(new)
}

// Flatten joins a 2-D slice of ints into a 1-D slice
func (m AoAoI) Flatten() AoI {
	if m.IsErr() {
		return ErrAoI(m.err)
	}

	xs := make([]int, 0)
	for _, v := range m.just {
		xs = append(xs, v...)
	}

	return JustAoI(xs)
}

// Map applies a function to each element of a valid AoAoI (i.e. a 1-D slice)
// and returns a new AoAoI.  If the AoAoI is invalid or if any function
// returns an invalid AoI, Map returns an invalid AoAoI.
func (m AoAoI) Map(f func(s []int) AoI) AoAoI {
	if m.IsErr() {
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
	if m.IsErr() {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// ToStr applies a function that takes an int and returns an S.  If the AoAoI
// is invalid or if any function returns an invalid S, ToStr returns an
// invalid AoAoS.  Note: unlike Map, this is a deep conversion of individual
// elements of the 2-D slice of ints.
func (m AoAoI) ToStr(f func(x int) S) AoAoS {
	if m.IsErr() {
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

// Unbox returns the underlying 2-D slice of ints or error.
func (m AoAoI) Unbox() ([][]int, error) {
	if m.just == nil && m.err == nil {
		return nil, errors.New("zero-value AoAoI")
	}
	return m.just, m.err
}
