// Package maybe implements the Maybe monad for some basic types plus arrays
// and 2-D arrays of those types.
//
// To keep type names short and manageable, abbreviations are used.  Type
// `maybe.I` is for ints; `maybe.AoI` is short for "array of ints" and
// `maybe.AoAoI` is short for "array of array of ints".
//
// This package only implements up to 2-D containers because those are common
// when working with line-oriented data.  For example, a text file can be
// interpreted as an array of an array of characters.
//
// Three constructors are provided for each type.  The `Just_` and `Err_`
// constructors are for values and errors, respectively.  The `New_`
// constructor can construct either type, and is intended for wrapping
// functions that follow the pattern of returning a value and an error.
package maybe
