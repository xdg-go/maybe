package maybe_test

import (
	"errors"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func TestString(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.S
	var just string
	var err error

	good = maybe.JustS("Hello")
	just, err = good.Unbox()
	is.Equal(just, "Hello")
	is.Nil(err)
	is.False(good.IsErr())

	bad = maybe.ErrS(errors.New("bad string"))
	just, err = bad.Unbox()
	is.Equal(just, "")
	is.NotNil(err)
	is.Equal(err.Error(), "bad string")
	is.True(bad.IsErr())

	got = maybe.NewS("Hello", nil)
	is.Equal(got, good)

	got = maybe.NewS("", err)
	is.True(got.IsErr())
}

func TestStringBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.S
	var just string
	var err error

	good = maybe.JustS("Hello")
	bad = maybe.ErrS(errors.New("bad string"))

	f := func(s string) maybe.S { return maybe.JustS(s + " World") }

	// Bind S to S; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, "Hello World")
	is.Nil(err)

	// Bind S to S; good path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestStringJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad maybe.S
	var got maybe.AoS
	var err error

	good = maybe.JustS("Hello")
	bad = maybe.ErrS(errors.New("bad string"))

	f := func(s string) maybe.AoS { return maybe.JustAoS([]string{s}) }

	// Split S to AoS
	got = good.Split(f)
	aos, err := got.Unbox()
	is.Equal(aos, []string{"Hello"})
	is.Nil(err)

	// Split S to AoS
	got = bad.Split(f)
	is.True(got.IsErr())
}
