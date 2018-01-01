package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getIntFixtures(input []int) (good, bad maybe.AoI) {
	good = maybe.JustAoI(input)
	bad = maybe.ErrAoI(errors.New("bad int"))
	return
}

func TestArrayOfInt(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []int{23, 42}
	good, bad := getIntFixtures(input)
	var got maybe.AoI
	var just []int
	var err error

	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)
	is.False(good.IsErr())

	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad int")
	is.True(bad.IsErr())

	got = maybe.NewAoI(input, nil)
	is.Equal(got, good)

	got = maybe.NewAoI(nil, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just [23 42]")
	is.Equal(bad.String(), "Err bad int")
}

func TestArrayOfIntBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []int{23, 42}
	good, bad := getIntFixtures(input)
	var got maybe.AoI
	var just []int
	var err error

	f := func(x []int) maybe.AoI { return maybe.JustAoI(x[1:]) }

	// Bind AoI to AoI; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, []int{42})
	is.Nil(err)

	// Bind AoI to AoI; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestArrayOfIntJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []int{23, 42}
	good, bad := getIntFixtures(input)
	var got maybe.I
	var err error

	f := func(xs []int) maybe.I {
		var sum int
		for _, v := range xs {
			sum += v
		}
		return maybe.JustI(sum)
	}

	// Join AoI to I; good path
	got = good.Join(f)
	x, err := got.Unbox()
	is.Equal(x, 65)
	is.Nil(err)

	// Join AoI to I; bad path
	got = bad.Join(f)
	is.True(got.IsErr())
}

func TestArrayOfIntMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []int{23, 42}
	good, bad := getIntFixtures(input)
	var just []int
	var err error

	// Map where everything succeeds
	neg := good.Map(func(x int) maybe.I { return maybe.JustI(-x) })
	just, err = neg.Unbox()
	is.Equal(just, []int{-23, -42})
	is.Nil(err)

	// Map where input is invalid
	negBadInput := bad.Map(func(x int) maybe.I { return maybe.JustI(-x) })
	is.True(negBadInput.IsErr())

	// Map where function returns invalid
	negBadMap := good.Map(func(x int) maybe.I { return maybe.ErrI(errors.New("bad int")) })
	is.True(negBadMap.IsErr())
}
