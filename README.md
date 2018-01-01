[![GoDoc](https://godoc.org/github.com/xdg/maybe?status.svg)](https://godoc.org/github.com/xdg/maybe)
[![Build Status](https://travis-ci.org/xdg/maybe.svg?branch=master)](https://travis-ci.org/xdg/maybe)

# maybe – A Maybe monad experiment for Go

## Synopsis

```go
import (
	"fmt"

	"github.com/xdg/maybe"
)

func Example() {

	// Various integers.
	nums := []int{42, 23, 0}

	// Wrap numbers in maybe.AoI ("Array of Integer") to "box" them.
	maybeNums := maybe.JustAoI(nums)

	// Define a function that, given a number, validates that it is positive
	// and returns either just the number or an error as a maybe.I.
	f := func(x int) maybe.I {
		if x < 0 {
			return maybe.ErrI(fmt.Errorf("%d is negative", x))
		}
		return maybe.JustI(x)
	}

	// Map that function onto boxed array of numbers.
	maybeNums = maybeNums.Map(f)

	// Examine the result
	if maybeNums.IsErr() {
		fmt.Println("Some value is negative")
	}

	fmt.Println("All values are zero or positive")
}
```

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

## Copyright and License

Copyright 2017 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
