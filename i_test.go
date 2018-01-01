package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func TestMaybeInt(t *testing.T) {
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

	// Bind I to I
	got = good.Bind(func(x int) maybe.I { return maybe.JustI(-x) })
	just, err = got.Unbox()
	is.Equal(just, -42)
	is.Nil(err)

	// Split I to AoI
	maos := good.Split(func(x int) maybe.AoI { return maybe.JustAoI([]int{x / 10, x % 10}) })
	aos, err := maos.Unbox()
	is.Equal(aos, []int{4, 2})
	is.Nil(err)
}
