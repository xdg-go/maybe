package maybe_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getFixtures(input []string) (good, bad maybe.AoS) {
	good = maybe.JustAoS(input)
	bad = maybe.ErrAoS(errors.New("bad string"))
	return
}

func TestMaybeArrayOfString(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"Hello", "World"}
	good, bad := getFixtures(input)
	var got maybe.AoS
	var just []string
	var err error

	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)
	is.False(good.IsErr())

	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad string")
	is.True(bad.IsErr())

	// Bind AoS to AoS
	got = good.Bind(func(s []string) maybe.AoS { return maybe.JustAoS(s[1:]) })
	just, err = got.Unbox()
	is.Equal(just, []string{"World"})
	is.Nil(err)

	// Join AoS to S
	ms := good.Join(func(s []string) maybe.S { return maybe.JustS(strings.Join(s, " ")) })
	s, err := ms.Unbox()
	is.Equal(s, "Hello World")
	is.Nil(err)
}

func TestMaybeArrayOfStringMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"Hello", "World"}
	good, bad := getFixtures(input)
	var just []string
	var err error

	// Map where everything succeeds
	lc := good.Map(func(s string) maybe.S { return maybe.JustS(strings.ToLower(s)) })
	just, err = lc.Unbox()
	is.Equal(just, []string{"hello", "world"})
	is.Nil(err)

	// Map where input is invalid
	lcBadInput := bad.Map(func(s string) maybe.S { return maybe.JustS(strings.ToLower(s)) })
	is.True(lcBadInput.IsErr())

	// Map where function returns invalid
	lcBadMap := good.Map(func(s string) maybe.S { return maybe.ErrS(errors.New("bad string")) })
	is.True(lcBadMap.IsErr())
}
