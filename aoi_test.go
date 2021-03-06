package maybe_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getIntFixtures(input []int) (good, bad maybe.AoI) {
	good = maybe.JustAoI(input)
	bad = maybe.ErrAoI(errors.New("bad int"))
	return
}

func TestAoI(t *testing.T) {
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

func TestAoIZero(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	// Check zero value case
	zero := maybe.AoI{}
	is.True(zero.IsErr())
	zero.Bind(func(x []int) maybe.AoI {
		if x == nil {
			panic("nil slice")
		}
		return maybe.JustAoI(x)
	})
	is.True(zero.IsErr())
	_, err := zero.Unbox()
	is.NotNil(err)
}

func TestAoIBind(t *testing.T) {
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

func TestAoISplit(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var got maybe.AoAoI
	var err error

	input := []int{23, 42}
	good, bad := getIntFixtures(input)

	f := func(x int) maybe.AoI { return maybe.JustAoI([]int{x}) }

	// Split S to AoI
	got = good.Split(f)
	aoaoi, err := got.Unbox()
	is.Equal(aoaoi, [][]int{[]int{23}, []int{42}})
	is.Nil(err)

	// Split S to AoI
	got = bad.Split(f)
	is.True(got.IsErr())

	// Split where input is invalid
	badSplit := good.Split(func(x int) maybe.AoI { return maybe.ErrAoI(errors.New("bad int")) })
	is.True(badSplit.IsErr())
}

func TestAoIJoin(t *testing.T) {
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

func TestAoIMap(t *testing.T) {
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

func TestAoIToString(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var got maybe.AoS
	var err error

	input := []int{23, 42}
	good, bad := getIntFixtures(input)

	f := func(x int) maybe.S { return maybe.JustS(fmt.Sprintf("%d", x)) }

	// Convert AoI to AoS; good path
	got = good.ToStr(f)
	s, err := got.Unbox()
	is.Equal(s, []string{"23", "42"})
	is.Nil(err)

	// Convert AoS to AoI; bad path
	got = bad.ToStr(f)
	is.True(got.IsErr())

	// Convert invalid I to S; albeit contrived
	got = good.ToStr(func(x int) maybe.S { return maybe.ErrS(errors.New("invalid")) })
	is.True(got.IsErr())
}
