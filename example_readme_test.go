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

// Example_StrsToInts shows how to convert a list of strings to a list of
// non-negative integers, accouting for the possibility of failure either in
// conversion or validation.
func Example_StrsToInts() {

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
