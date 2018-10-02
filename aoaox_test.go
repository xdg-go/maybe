package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getAoAoXFixtures(input [][]interface{}) (good, bad maybe.AoAoX) {
	good = maybe.JustAoAoX(input)
	bad = maybe.ErrAoAoX(errors.New("bad interface{}s"))
	return
}

func TestAoAoX(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]interface{}{
		[]interface{}{23, 42},
		[]interface{}{11, 13},
	}
	good, bad := getAoAoXFixtures(input)
	var got maybe.AoAoX
	var just [][]interface{}
	var err error

	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)
	is.False(good.IsErr())

	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad interface{}s")
	is.True(bad.IsErr())

	got = maybe.NewAoAoX(input, nil)
	is.Equal(got, good)

	got = maybe.NewAoAoX(nil, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just [[23 42] [11 13]]")
	is.Equal(bad.String(), "Err bad interface{}s")
}

func TestAoAoXFromSlice(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]interface{}{
		[]interface{}{23, 42},
		[]interface{}{11, 13},
	}
	good := maybe.JustAoAoX(input)
	var got maybe.AoAoX

	// Create from slice and non-slice
	goodInput := [][]int{
		[]int{23, 42},
		[]int{11, 13},
	}
	badInput1 := 42
	badInput2 := []int{23, 42}

	got = maybe.NewAoAoXFromSlice(goodInput, nil)
	is.False(got.IsErr())
	is.Equal(got, good)

	got = maybe.NewAoAoXFromSlice(badInput1, nil)
	is.True(got.IsErr())

	got = maybe.NewAoAoXFromSlice(badInput2, nil)
	is.True(got.IsErr())

	got = maybe.NewAoAoXFromSlice(nil, nil)
	is.True(got.IsErr())

	got = maybe.NewAoAoXFromSlice(nil, errors.New("from error"))
	is.True(got.IsErr())
	_, err := got.Unbox()
	is.Equal(err.Error(), "from error")
}

func TestAoAoXZero(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	// Check zero value case
	zero := maybe.AoAoX{}
	is.True(zero.IsErr())
	zero.Bind(func(x [][]interface{}) maybe.AoAoX {
		if x == nil {
			panic("nil slice")
		}
		return maybe.JustAoAoX(x)
	})
	is.True(zero.IsErr())
	_, err := zero.Unbox()
	is.NotNil(err)
}
func TestAoAoXBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]interface{}{
		[]interface{}{23, 42},
		[]interface{}{11, 13},
	}
	good, bad := getAoAoXFixtures(input)
	var got maybe.AoAoX
	var just [][]interface{}
	var err error

	f := func(s [][]interface{}) maybe.AoAoX { return maybe.JustAoAoX(s[1:]) }

	// Bind AoAoX to AoAoX; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, [][]interface{}{[]interface{}{11, 13}})
	is.Nil(err)

	// Bind AoAoX to AoAoX; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestAoAoXJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]interface{}{
		[]interface{}{23, 42},
		[]interface{}{11, 13},
	}
	good, bad := getAoAoXFixtures(input)
	var got maybe.AoX
	var err error

	f := func(x []interface{}) maybe.X {
		sum := 0
		for _, v := range x {
			sum += v.(int)
		}
		return maybe.JustX(sum)
	}

	// Join AoAoX to AoX; good path
	got = good.Join(f)
	s, err := got.Unbox()
	is.Equal(s, []interface{}{65, 24})
	is.Nil(err)

	// Join AoAoX to AoX; bad path
	got = bad.Join(f)
	is.True(got.IsErr())

	// Join where input is invalid
	badJoin := good.Join(func(x []interface{}) maybe.X { return maybe.ErrX(errors.New("bad interface{}")) })
	is.True(badJoin.IsErr())
}

func TestAoAoXMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]interface{}{
		[]interface{}{23, 42},
		[]interface{}{11, 13},
	}
	good, bad := getAoAoXFixtures(input)
	var just [][]interface{}
	var err error

	// Map where everything succeeds
	// Function returns first element only
	f := func(xs []interface{}) maybe.AoX { return maybe.JustAoX(xs[0:1]) }
	firsts := good.Map(f)
	just, err = firsts.Unbox()
	is.Equal(just, [][]interface{}{[]interface{}{23}, []interface{}{11}})
	is.Nil(err)

	// Map where input is invalid
	badInput := bad.Map(f)
	is.True(badInput.IsErr())

	// Map where function returns invalid
	badMap := good.Map(func(xs []interface{}) maybe.AoX { return maybe.ErrAoX(errors.New("bad interface{}s")) })
	is.True(badMap.IsErr())
}

func TestAoAoXFlatten(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]interface{}{
		[]interface{}{23, 42},
		[]interface{}{11, 13},
	}
	good, bad := getAoAoXFixtures(input)
	var got maybe.AoX
	var just []interface{}
	var err error

	// Good path
	got = good.Flatten()
	just, err = got.Unbox()
	is.Equal(just, []interface{}{23, 42, 11, 13})
	is.Nil(err)

	// Bad path
	got = bad.Flatten()
	is.True(got.IsErr())
}
