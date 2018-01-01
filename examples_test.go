package maybe_test

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
