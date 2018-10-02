package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getAoXFixtures(input []interface{}) (good, bad maybe.AoX) {
	good = maybe.JustAoX(input)
	bad = maybe.ErrAoX(errors.New("bad interface{}"))
	return
}

func TestAoX(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []interface{}{23, 42}
	good, bad := getAoXFixtures(input)
	var got maybe.AoX
	var just []interface{}
	var err error

	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)
	is.False(good.IsErr())

	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad interface{}")
	is.True(bad.IsErr())

	got = maybe.NewAoX(input, nil)
	is.Equal(got, good)

	got = maybe.NewAoX(nil, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just [23 42]")
	is.Equal(bad.String(), "Err bad interface{}")
}

func TestAoXFromSlice(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []interface{}{23, 42}
	good := maybe.JustAoX(input)
	var got maybe.AoX

	// Create from slice and non-slice
	goodInput := []int{23, 42}
	badInput := 42

	got = maybe.NewAoXFromSlice(goodInput, nil)
	is.False(got.IsErr())
	is.Equal(got, good)

	got = maybe.NewAoXFromSlice(badInput, nil)
	is.True(got.IsErr())

	got = maybe.NewAoXFromSlice(nil, nil)
	is.True(got.IsErr())

	got = maybe.NewAoXFromSlice(nil, errors.New("from error"))
	is.True(got.IsErr())
	_, err := got.Unbox()
	is.Equal(err.Error(), "from error")
}

func TestAoXZero(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	// Check zero value case
	zero := maybe.AoX{}
	is.True(zero.IsErr())
	zero.Bind(func(x []interface{}) maybe.AoX {
		if x == nil {
			panic("nil slice")
		}
		return maybe.JustAoX(x)
	})
	is.True(zero.IsErr())
	_, err := zero.Unbox()
	is.NotNil(err)
}

func TestAoXBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []interface{}{23, 42}
	good, bad := getAoXFixtures(input)
	var got maybe.AoX
	var just []interface{}
	var err error

	f := func(x []interface{}) maybe.AoX { return maybe.JustAoX(x[1:]) }

	// Bind AoX to AoX; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, []interface{}{42})
	is.Nil(err)

	// Bind AoX to AoX; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestAoXSplit(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var got maybe.AoAoX
	var err error

	input := []interface{}{23, 42}
	good, bad := getAoXFixtures(input)

	f := func(x interface{}) maybe.AoX { return maybe.JustAoX([]interface{}{x}) }

	// Split S to AoX
	got = good.Split(f)
	aoaoi, err := got.Unbox()
	is.Equal(aoaoi, [][]interface{}{[]interface{}{23}, []interface{}{42}})
	is.Nil(err)

	// Split S to AoX
	got = bad.Split(f)
	is.True(got.IsErr())

	// Split where input is invalid
	badSplit := good.Split(func(x interface{}) maybe.AoX { return maybe.ErrAoX(errors.New("bad interface{}")) })
	is.True(badSplit.IsErr())
}

func TestAoXJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []interface{}{23, 42}
	good, bad := getAoXFixtures(input)
	var got maybe.X
	var err error

	f := func(xs []interface{}) maybe.X {
		var sum int
		for _, v := range xs {
			sum += v.(int)
		}
		return maybe.JustX(sum)
	}

	// Join AoX to I; good path
	got = good.Join(f)
	x, err := got.Unbox()
	is.Equal(x, 65)
	is.Nil(err)

	// Join AoX to I; bad path
	got = bad.Join(f)
	is.True(got.IsErr())
}

func TestAoXMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []interface{}{23, 42}
	good, bad := getAoXFixtures(input)
	var just []interface{}
	var err error

	// Map where everything succeeds
	neg := good.Map(func(x interface{}) maybe.X { return maybe.JustX(-x.(int)) })
	just, err = neg.Unbox()
	is.Equal(just, []interface{}{-23, -42})
	is.Nil(err)

	// Map where input is invalid
	negBadInput := bad.Map(func(x interface{}) maybe.X { return maybe.JustX(-x.(int)) })
	is.True(negBadInput.IsErr())

	// Map where function returns invalid
	negBadMap := good.Map(func(x interface{}) maybe.X { return maybe.ErrX(errors.New("bad interface{}")) })
	is.True(negBadMap.IsErr())
}
