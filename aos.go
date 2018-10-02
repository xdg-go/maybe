package maybe

import (
	"errors"
	"fmt"
)

// AoS implements the Maybe monad for a slice of strings.  An AoS is
// considered 'valid' or 'invalid' depending on whether it contains a slice of
// strings or an error value.  A zero-value AoS is invalid and Unbox() will
// return an error to that effect.
type AoS struct {
	just []string
	err  error
}

// NewAoS constructs an AoS from a given slice of strings or error. If e is
// not nil, returns ErrAoS(e), otherwise returns JustAoS(s)
func NewAoS(s []string, e error) AoS {
	if e != nil {
		return ErrAoS(e)
	}
	return JustAoS(s)
}

// JustAoS constructs a valid AoS from a given slice of strings.
func JustAoS(s []string) AoS {
	return AoS{just: s}
}

// ErrAoS constructs an invalid AoS from a given error.
func ErrAoS(e error) AoS {
	return AoS{err: e}
}

// IsErr returns true for an invalid AoS.
func (m AoS) IsErr() bool {
	return m.just == nil || m.err != nil
}

// Bind applies a function that takes a slice of strings and returns an AoS.
func (m AoS) Bind(f func(s []string) AoS) AoS {
	if m.IsErr() {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a slice of strings and returns an S.
func (m AoS) Join(f func(s []string) S) S {
	if m.IsErr() {
		return ErrS(m.err)
	}

	return f(m.just)
}

// Split applies a splitting function to each element of a valid AoS,
// resulting in a higher-dimension structure. If the AoS is invalid or if any
// function returns an invalid AoS, Split returns an invalid AoAoS.
func (m AoS) Split(f func(s string) AoS) AoAoS {
	if m.IsErr() {
		return ErrAoAoS(m.err)
	}

	new := make([][]string, len(m.just))
	for i, v := range m.just {
		xs, err := f(v).Unbox()
		if err != nil {
			return ErrAoAoS(err)
		}
		new[i] = xs
	}

	return JustAoAoS(new)
}

// Map applies a function to each element of a valid AoS and returns a new
// AoS.  If the AoS is invalid or if any function returns an invalid S, Map
// returns an invalid AoS.
func (m AoS) Map(f func(s string) S) AoS {
	if m.IsErr() {
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

// ToInt applies a function that takes a string and returns an I.If the AoS is
// invalid or if any function returns an invalid I, ToInt returns an invalid
// AoI.
func (m AoS) ToInt(f func(s string) I) AoI {
	if m.IsErr() {
		return ErrAoI(m.err)
	}

	new := make([]int, len(m.just))
	for i, v := range m.just {
		num, err := f(v).Unbox()
		if err != nil {
			return ErrAoI(err)
		}
		new[i] = num
	}

	return JustAoI(new)
}

// String returns a string representation, mostly useful for debugging.
func (m AoS) String() string {
	if m.IsErr() {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// Unbox returns the underlying slice of strings value or error.
func (m AoS) Unbox() ([]string, error) {
	if m.just == nil && m.err == nil {
		return nil, errors.New("zero-value AoS")
	}
	return m.just, m.err
}
