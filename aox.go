package maybe

import (
	"errors"
	"fmt"
	"reflect"
)

// AoX implements the Maybe monad for a slice of empty interfaces.  An AoX is
// considered 'valid' or 'invalid' depending on whether it contains a slice of
// empty interfaces or an error value.  A zero-value AoX is invalid and
// Unbox() will return an error to that effect.
type AoX struct {
	just []interface{}
	err  error
}

// NewAoX constructs an AoX from a given slice of empty interfaces or error.
// If e is not nil, returns ErrAoX(e), otherwise returns JustAoX(x).
func NewAoX(x []interface{}, e error) AoX {
	if e != nil {
		return ErrAoX(e)
	}
	return JustAoX(x)
}

var errAoXNotSlice = errors.New("NewAoXFromSlice called with non-slice")

// NewAoXFromSlice constructs an AoX from a given slice of arbitrary values or
// error.  If e is not nil, returns ErrAoX(e), otherwise, the slice of values
// is converted to a slice of empty interface and returned as JustAoX(x).  If
// the provided value is not a slice, ErrAoX is returned.
func NewAoXFromSlice(x interface{}, e error) AoX {
	if e != nil {
		return ErrAoX(e)
	}
	if x == nil {
		return ErrAoX(errAoXNotSlice)
	}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(x)
		xs := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			xs[i] = s.Index(i).Interface()
		}
		return JustAoX(xs)
	default:
		return ErrAoX(errAoXNotSlice)
	}
}

// JustAoX constructs a valid AoX from a given slice of empty interfaces.
func JustAoX(x []interface{}) AoX {
	return AoX{just: x}
}

// ErrAoX constructs an invalid AoX from a given error.
func ErrAoX(e error) AoX {
	return AoX{err: e}
}

// IsErr returns true for an invalid AoX.
func (m AoX) IsErr() bool {
	return m.just == nil || m.err != nil
}

// Bind applies a function that takes a slice of empty interfaces and returns
// an AoX.
func (m AoX) Bind(f func(x []interface{}) AoX) AoX {
	if m.IsErr() {
		return m
	}

	return f(m.just)
}

// Join applies a function that takes a slice of empty interfaces and returns
// an I.
func (m AoX) Join(f func(x []interface{}) X) X {
	if m.IsErr() {
		return ErrX(m.err)
	}

	return f(m.just)
}

// Split applies a splitting function to each element of a valid AoX,
// resulting in a higher-dimension structure. If the AoX is invalid or if any
// function returns an invalid AoX, Split returns an invalid AoAoX.
func (m AoX) Split(f func(x interface{}) AoX) AoAoX {
	if m.IsErr() {
		return ErrAoAoX(m.err)
	}

	xss := make([][]interface{}, len(m.just))
	for i, v := range m.just {
		xs, err := f(v).Unbox()
		if err != nil {
			return ErrAoAoX(err)
		}
		xss[i] = xs
	}

	return JustAoAoX(xss)
}

// Map applies a function to each element of a valid AoX and returns a new
// AoX.  If the AoX is invalid or if any function returns an invalid I, Map
// returns an invalid AoX.
func (m AoX) Map(f func(x interface{}) X) AoX {
	if m.IsErr() {
		return m
	}

	xss := make([]interface{}, len(m.just))
	for i, v := range m.just {
		x, err := f(v).Unbox()
		if err != nil {
			return ErrAoX(err)
		}
		xss[i] = x
	}

	return JustAoX(xss)
}

// String returns a string representation, mostly useful for debugging.
func (m AoX) String() string {
	if m.IsErr() {
		return fmt.Sprintf("Err %v", m.err)
	}
	return fmt.Sprintf("Just %v", m.just)
}

// Unbox returns the underlying slice of empty interfaces or error.
func (m AoX) Unbox() ([]interface{}, error) {
	if m.just == nil && m.err == nil {
		return nil, errors.New("zero-value AoX")
	}
	return m.just, m.err
}
