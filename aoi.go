package maybe

import (
	"errors"
	"fmt"
)

// AoI implements the Maybe monad for a slice of ints.  An AoI is considered
// 'valid' or 'invalid' depending on whether it contains a slice of ints or an
// error value.  A zero-value AoI is invalid and Unbox() will return an error
// to that effect.
type AoI struct {
	just []int
	err  error
}

// NewAoI constructs an AoI from a given slice of ints or error. If e is not
// nil, returns ErrAoI(e), otherwise returns JustAoI(s).
func NewAoI(s []int, e error) AoI {
	if e != nil {
		return ErrAoI(e)
	}
	return JustAoI(s)
}

// JustAoI constructs a valid AoI from a given slice of ints.
func JustAoI(s []int) AoI {
	return AoI{just: s}
}

// ErrAoI constructs an invalid AoI from a given error.
func ErrAoI(e error) AoI {
	return AoI{err: e}
}

// IsErr returns true for an invalid AoI.
func (m AoI) IsErr() bool {
	return m.just == nil || m.err != nil
}

// Bind applies a function that takes a slice of ints and returns an AoI.
func (m AoI) Bind(f func(s []int) AoI) AoI {
	if m.IsErr() {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a slice of ints and returns an I.
func (m AoI) Join(f func(s []int) I) I {
	if m.IsErr() {
		return ErrI(m.err)
	}

	return f(m.just)
}

// Map applies a function to each element of a valid AoI and returns a new
// AoI.  If the AoI is invalid or if any function returns an invalid I, Map
// returns an invalid AoI.
func (m AoI) Map(f func(s int) I) AoI {
	if m.IsErr() {
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

// String returns a string representation, mostly useful for debugging.
func (m AoI) String() string {
	if m.IsErr() {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// ToStr applies a function that takes an int and returns an S.  If the AoI
// is invalid or if any function returns an invalid S, ToStr returns an
// invalid AoS.
func (m AoI) ToStr(f func(x int) S) AoS {
	if m.IsErr() {
		return ErrAoS(m.err)
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

// Unbox returns the underlying slice of ints or error.
func (m AoI) Unbox() ([]int, error) {
	if m.just == nil && m.err == nil {
		return nil, errors.New("zero-value AoI")
	}
	return m.just, m.err
}
