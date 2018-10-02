package maybe_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getAoAoIFixtures(input [][]int) (good, bad maybe.AoAoI) {
	good = maybe.JustAoAoI(input)
	bad = maybe.ErrAoAoI(errors.New("bad ints"))
	return
}

func TestAoAoI(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]int{
		[]int{23, 42},
		[]int{11, 13},
	}
	good, bad := getAoAoIFixtures(input)
	var got maybe.AoAoI
	var just [][]int
	var err error

	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)
	is.False(good.IsErr())

	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad ints")
	is.True(bad.IsErr())

	got = maybe.NewAoAoI(input, nil)
	is.Equal(got, good)

	got = maybe.NewAoAoI(nil, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just [[23 42] [11 13]]")
	is.Equal(bad.String(), "Err bad ints")
}

func TestAoAoIZero(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	// Check zero value case
	zero := maybe.AoAoI{}
	is.True(zero.IsErr())
	zero.Bind(func(x [][]int) maybe.AoAoI {
		if x == nil {
			panic("nil slice")
		}
		return maybe.JustAoAoI(x)
	})
	is.True(zero.IsErr())
	_, err := zero.Unbox()
	is.NotNil(err)
}
func TestAoAoIBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]int{
		[]int{23, 42},
		[]int{11, 13},
	}
	good, bad := getAoAoIFixtures(input)
	var got maybe.AoAoI
	var just [][]int
	var err error

	f := func(s [][]int) maybe.AoAoI { return maybe.JustAoAoI(s[1:]) }

	// Bind AoAoI to AoAoI; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, [][]int{[]int{11, 13}})
	is.Nil(err)

	// Bind AoAoI to AoAoI; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestAoAoIJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]int{
		[]int{23, 42},
		[]int{11, 13},
	}
	good, bad := getAoAoIFixtures(input)
	var got maybe.AoI
	var err error

	f := func(x []int) maybe.I {
		sum := 0
		for _, v := range x {
			sum += v
		}
		return maybe.JustI(sum)
	}

	// Join AoAoI to AoI; good path
	got = good.Join(f)
	s, err := got.Unbox()
	is.Equal(s, []int{65, 24})
	is.Nil(err)

	// Join AoAoI to AoI; bad path
	got = bad.Join(f)
	is.True(got.IsErr())
}

func TestAoAoIMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]int{
		[]int{23, 42},
		[]int{11, 13},
	}
	good, bad := getAoAoIFixtures(input)
	var just [][]int
	var err error

	// Map where everything succeeds
	// Function returns first element only
	f := func(xs []int) maybe.AoI { return maybe.JustAoI(xs[0:1]) }
	firsts := good.Map(f)
	just, err = firsts.Unbox()
	is.Equal(just, [][]int{[]int{23}, []int{11}})
	is.Nil(err)

	// Map where input is invalid
	badInput := bad.Map(f)
	is.True(badInput.IsErr())

	// Map where function returns invalid
	badMap := good.Map(func(xs []int) maybe.AoI { return maybe.ErrAoI(errors.New("bad ints")) })
	is.True(badMap.IsErr())
}

func TestAoAoIToStr(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]int{
		[]int{23, 42},
		[]int{11, 13},
	}
	good, bad := getAoAoIFixtures(input)
	var got maybe.AoAoS
	var err error

	f := func(n int) maybe.S { return maybe.JustS(strconv.Itoa(n)) }

	// Convert I to S; good path
	got = good.ToStr(f)
	x, err := got.Unbox()
	is.Equal(x, [][]string{[]string{"23", "42"}, []string{"11", "13"}})
	is.Nil(err)

	// Convert invalid I to S
	got = bad.ToStr(f)
	is.True(got.IsErr())

	// Convert invalid I to S; albeit contrived
	got = good.ToStr(func(x int) maybe.S { return maybe.ErrS(errors.New("invalid")) })
	is.True(got.IsErr())
}
