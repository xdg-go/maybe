package maybe

import (
	"errors"
	"fmt"
	"reflect"
)

// AoAoX implements the Maybe monad for a 2-D slice of empty interfaces.  An AoAoX is
// considered 'valid' or 'invalid' depending on whether it contains a 2-D
// slice of empty interfaces or an error value.  A zero-value AoAoX is invalid and Unbox()
// will return an error to that effect.
type AoAoX struct {
	just [][]interface{}
	err  error
}

// NewAoAoX constructs an AoAoX from a given 2-D slice of empty interfaces or error. If e is not
// nil, returns ErrAoAoX(e), otherwise returns JustAoAoX(x).
func NewAoAoX(x [][]interface{}, e error) AoAoX {
	if e != nil {
		return ErrAoAoX(e)
	}
	return JustAoAoX(x)
}

var errAoAoXNotSlice = errors.New("NewAoAoXFromSlice called with non-slice-of-slices")

// NewAoAoXFromSlice constructs an AoAoX from a given slice of slices of
// arbitrary values or error.  If e is not nil, returns ErrAoAoX(e),
// otherwise, the inner slices of values are converted to slices of empty
// interface and returned as JustAoAoX(x).  If the provided value is not a
// slice of slices, ErrAoAoX is returned.
func NewAoAoXFromSlice(x interface{}, e error) AoAoX {
	if e != nil {
		return ErrAoAoX(e)
	}
	if x == nil {
		return ErrAoAoX(errAoAoXNotSlice)
	}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(x)
		xs := make([][]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			v, err := NewAoXFromSlice(s.Index(i).Interface(), nil).Unbox()
			if err != nil {
				return ErrAoAoX(errAoAoXNotSlice)
			}
			xs[i] = v
		}
		return JustAoAoX(xs)
	default:
		return ErrAoAoX(errAoAoXNotSlice)
	}
}

// JustAoAoX constructs a valid AoAoX from a given 2-D slice of empty interfaces.
func JustAoAoX(x [][]interface{}) AoAoX {
	return AoAoX{just: x}
}

// ErrAoAoX constructs an invalid AoAoX from a given error.
func ErrAoAoX(e error) AoAoX {
	return AoAoX{err: e}
}

// IsErr returns true for an invalid AoAoX.
func (m AoAoX) IsErr() bool {
	return m.just == nil || m.err != nil
}

// Bind applies a function that takes a 2-D slice of empty interfaces and returns an AoAoX.
func (m AoAoX) Bind(f func(x [][]interface{}) AoAoX) AoAoX {
	if m.IsErr() {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a 2-D slice of empty interfaces and returns an AoX.
func (m AoAoX) Join(f func(x []interface{}) X) AoX {
	if m.IsErr() {
		return ErrAoX(m.err)
	}

	xss := make([]interface{}, len(m.just))
	for i, v := range m.just {
		s, err := f(v).Unbox()
		if err != nil {
			return ErrAoX(err)
		}
		xss[i] = s
	}

	return JustAoX(xss)
}

// Flatten joins a 2-D slice of empty interfaces into a 1-D slice
func (m AoAoX) Flatten() AoX {
	if m.IsErr() {
		return ErrAoX(m.err)
	}

	xs := make([]interface{}, 0)
	for _, v := range m.just {
		xs = append(xs, v...)
	}

	return JustAoX(xs)
}

// Map applies a function to each element of a valid AoAoX (i.e. a 1-D slice)
// and returns a new AoAoX.  If the AoAoX is invalid or if any function
// returns an invalid AoX, Map returns an invalid AoAoX.
func (m AoAoX) Map(f func(x []interface{}) AoX) AoAoX {
	if m.IsErr() {
		return m
	}

	xss := make([][]interface{}, len(m.just))
	for i, v := range m.just {
		x, err := f(v).Unbox()
		if err != nil {
			return ErrAoAoX(err)
		}
		xss[i] = x
	}

	return JustAoAoX(xss)
}

// String returns a string representation, mostly useful for debugging.
func (m AoAoX) String() string {
	if m.IsErr() {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// Unbox returns the underlying 2-D slice of empty interfaces or error.
func (m AoAoX) Unbox() ([][]interface{}, error) {
	if m.just == nil && m.err == nil {
		return nil, errors.New("zero-value AoAoX")
	}
	return m.just, m.err
}
