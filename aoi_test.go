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

func TestMaybeArrayOfInt(t *testing.T) {
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

	// Bind AoI to AoI
	got = good.Bind(func(x []int) maybe.AoI { return maybe.JustAoI(x[1:]) })
	just, err = got.Unbox()
	is.Equal(just, []int{42})
	is.Nil(err)

	// Join AoI to I
	ms := good.Join(func(xs []int) maybe.I {
		var sum int
		for _, v := range xs {
			sum += v
		}
		return maybe.JustI(sum)
	})
	x, err := ms.Unbox()
	is.Equal(x, 65)
	is.Nil(err)
}

func TestMaybeArrayOfIntMap(t *testing.T) {
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
