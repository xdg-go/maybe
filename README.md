[![GoDoc](https://godoc.org/github.com/xdg/maybe?status.svg)](https://godoc.org/github.com/xdg/maybe)
[![Build Status](https://travis-ci.org/xdg/maybe.svg?branch=master)](https://travis-ci.org/xdg/maybe)

# maybe – A Maybe monad experiment for Go

## Description

This package implements the [Maybe
monad](https://en.wikipedia.org/wiki/Monad_(functional_programming)#The_Maybe_monad)
for a couple basic types and arrays of those types.  This allows "boxing" a
value or error and chaining operations on the boxed type without constant
error checking.  See ["Error Handling in Go: Rob Pike Reinvented
Monads"](https://www.innoq.com/en/blog/golang-errors-monads/) for more on
the concept.

This is an experiment to simplify some libraries the author is writing.  It
should not be considered stable for production use.

## Example

```go
package maybe_test

import (
	"fmt"
	"strconv"

	"github.com/xdg/maybe"
)

type example struct {
	label string
	data  []string
}

// Example shows how to convert a list of strings to a list of non-negative
// integers, accouting for the possibility of failure either in conversion or
// validation.
func Example() {

	cases := []example{
		{label: "success", data: []string{"23", "42", "0"}},
		{label: "bad atoi", data: []string{"23", "forty-two", "0"}},
		{label: "negative", data: []string{"23", "-42", "0"}},
	}

	// Function to convert string to maybe.I.
	atoi := func(s string) maybe.I { return maybe.NewI(strconv.Atoi(s)) }

	// Function to validate non-negative integer.
	validate := func(x int) maybe.I {
		if x < 0 {
			return maybe.ErrI(fmt.Errorf("%d is negative", x))
		}
		return maybe.JustI(x)
	}

	// For each example, try converting and validating functionally and
	// then inspecting the result.
	for _, c := range cases {
		// Wrap the []string in a maybe type.
		strs := maybe.JustAoS(c.data)

		// Functionally convert and validate.
		nums := strs.ToInt(atoi).Map(validate)

		// Check if it worked.
		if nums.IsErr() {
			fmt.Printf("%s: %v failed to convert: %v\n", c.label, strs, nums)
		} else {
			fmt.Printf("%s: %v converted to %v\n", c.label, strs, nums)
		}
	}

	// Output:
	// success: Just [23 42 0] converted to Just [23 42 0]
	// bad atoi: Just [23 forty-two 0] failed to convert: Err strconv.Atoi: parsing "forty-two": invalid syntax
	// negative: Just [23 -42 0] failed to convert: Err -42 is negative
}
```

## Copyright and License

Copyright 2017 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
