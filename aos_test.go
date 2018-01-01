package maybe_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getStrFixtures(input []string) (good, bad maybe.AoS) {
	good = maybe.JustAoS(input)
	bad = maybe.ErrAoS(errors.New("bad string"))
	return
}

func TestArrayOfString(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"Hello", "World"}
	good, bad := getStrFixtures(input)
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

	got = maybe.NewAoS(input, nil)
	is.Equal(got, good)

	got = maybe.NewAoS(nil, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just [Hello World]")
	is.Equal(bad.String(), "Err bad string")
}

func TestArrayOfStringBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"Hello", "World"}
	good, bad := getStrFixtures(input)
	var got maybe.AoS
	var just []string
	var err error

	f := func(s []string) maybe.AoS { return maybe.JustAoS(s[1:]) }

	// Bind AoS to AoS; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, []string{"World"})
	is.Nil(err)

	// Bind AoS to AoS; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestArrayOfStringJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"Hello", "World"}
	good, bad := getStrFixtures(input)
	var got maybe.S
	var err error

	f := func(s []string) maybe.S { return maybe.JustS(strings.Join(s, " ")) }

	// Join AoS to S; good path
	got = good.Join(f)
	s, err := got.Unbox()
	is.Equal(s, "Hello World")
	is.Nil(err)

	// Join AoS to S; bad path
	got = bad.Join(f)
	is.True(got.IsErr())
}

func TestArrayOfStringMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"Hello", "World"}
	good, bad := getStrFixtures(input)
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

func TestArrayOfStringToInt(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := []string{"42", "23"}
	good, bad := getStrFixtures(input)
	notNum := maybe.JustAoS([]string{"23", "forty-two"})
	var got maybe.AoI
	var err error

	f := func(s string) maybe.I { return maybe.NewI(strconv.Atoi(s)) }

	// Convert S to I; good path
	got = good.ToInt(f)
	x, err := got.Unbox()
	is.Equal(x, []int{42, 23})
	is.Nil(err)

	// Convert S to I; bad path
	got = notNum.ToInt(f)
	is.True(got.IsErr())

	// Convert invalid S to I
	got = bad.ToInt(f)
	is.True(got.IsErr())
}
