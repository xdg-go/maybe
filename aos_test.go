package maybe_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func TestMaybeArrayOfString(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	var good, bad, got maybe.AoS
	var just []string
	var err error

	input := []string{"Hello", "World"}
	good = maybe.NewAoS(input)
	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)

	bad = maybe.ErrAoS(errors.New("bad string"))
	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad string")

	// Map AoS to AoS
	got = good.Bind(func(s []string) maybe.AoS { return maybe.NewAoS(s[1:]) })
	just, err = got.Unbox()
	is.Equal(just, []string{"World"})
	is.Nil(err)

	// Map AoS to S
	ms := good.BindS(func(s []string) maybe.S { return maybe.NewS(strings.Join(s, " ")) })
	s, err := ms.Unbox()
	is.Equal(s, "Hello World")
	is.Nil(err)
}
