package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func TestX(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.X
	var just interface{}
	var err error

	good = maybe.JustX(42)
	just, err = good.Unbox()
	is.Equal(just, 42)
	is.Nil(err)
	is.False(good.IsErr())

	bad = maybe.ErrX(errors.New("bad interface{}"))
	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad interface{}")
	is.True(bad.IsErr())

	got = maybe.NewX(42, nil)
	is.Equal(got, good)

	got = maybe.NewX(0, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just 42")
	is.Equal(bad.String(), "Err bad interface{}")
}

func TestXBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.X
	var just interface{}
	var err error

	good = maybe.JustX(42)
	bad = maybe.ErrX(errors.New("bad interface{}"))

	f := func(x interface{}) maybe.X { return maybe.JustX(-(x.(int))) }

	// Bind I to I; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, -42)
	is.Nil(err)

	// Bind I to I; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestXSplit(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad maybe.X
	var got maybe.AoX
	var err error

	good = maybe.JustX(42)
	bad = maybe.ErrX(errors.New("bad interface{}"))

	f := func(x interface{}) maybe.AoX { return maybe.JustAoX([]interface{}{x.(int) / 10, x.(int) % 10}) }

	// Split I to AoX
	got = good.Split(f)
	aos, err := got.Unbox()
	is.Equal(aos, []interface{}{4, 2})
	is.Nil(err)

	// Split I to AoX
	got = bad.Split(f)
	is.True(got.IsErr())
}
