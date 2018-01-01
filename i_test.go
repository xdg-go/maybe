package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func TestInt(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.I
	var just int
	var err error

	good = maybe.JustI(42)
	just, err = good.Unbox()
	is.Equal(just, 42)
	is.Nil(err)
	is.False(good.IsErr())

	bad = maybe.ErrI(errors.New("bad int"))
	just, err = bad.Unbox()
	is.Equal(just, 0)
	is.NotNil(err)
	is.Equal(err.Error(), "bad int")
	is.True(bad.IsErr())

	got = maybe.NewI(42, nil)
	is.Equal(got, good)

	got = maybe.NewI(0, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just 42")
	is.Equal(bad.String(), "Err bad int")
}

func TestIntBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.I
	var just int
	var err error

	good = maybe.JustI(42)
	bad = maybe.ErrI(errors.New("bad int"))

	f := func(x int) maybe.I { return maybe.JustI(-x) }

	// Bind I to I; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, -42)
	is.Nil(err)

	// Bind I to I; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestIntJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad maybe.I
	var got maybe.AoI
	var err error

	good = maybe.JustI(42)
	bad = maybe.ErrI(errors.New("bad int"))

	f := func(x int) maybe.AoI { return maybe.JustAoI([]int{x / 10, x % 10}) }

	// Split I to AoI
	got = good.Split(f)
	aos, err := got.Unbox()
	is.Equal(aos, []int{4, 2})
	is.Nil(err)

	// Split I to AoI
	got = bad.Split(f)
	is.True(got.IsErr())
}
