package maybe_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/xdg/maybe"
	"github.com/xdg/testy"
)

func getAoAoSFixtures(input [][]string) (good, bad maybe.AoAoS) {
	good = maybe.JustAoAoS(input)
	bad = maybe.ErrAoAoS(errors.New("bad strings"))
	return
}

func TestAoAoS(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]string{
		[]string{"Hello", "World"},
		[]string{"Goodbye", "Cruel World"},
	}
	good, bad := getAoAoSFixtures(input)
	var got maybe.AoAoS
	var just [][]string
	var err error

	just, err = good.Unbox()
	is.Equal(just, input)
	is.Nil(err)
	is.False(good.IsErr())

	just, err = bad.Unbox()
	is.Nil(just)
	is.NotNil(err)
	is.Equal(err.Error(), "bad strings")
	is.True(bad.IsErr())

	got = maybe.NewAoAoS(input, nil)
	is.Equal(got, good)

	got = maybe.NewAoAoS(nil, err)
	is.True(got.IsErr())

	is.Equal(good.String(), "Just [[Hello World] [Goodbye Cruel World]]")
	is.Equal(bad.String(), "Err bad strings")
}

func TestAoAoSZero(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	// Check zero value case
	zero := maybe.AoAoS{}
	is.True(zero.IsErr())
	zero.Bind(func(x [][]string) maybe.AoAoS {
		if x == nil {
			panic("nil slice")
		}
		return maybe.JustAoAoS(x)
	})
	is.True(zero.IsErr())
	_, err := zero.Unbox()
	is.NotNil(err)
}

func TestAoAoSBind(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]string{
		[]string{"Hello", "World"},
		[]string{"Goodbye", "Cruel World"},
	}
	good, bad := getAoAoSFixtures(input)
	var got maybe.AoAoS
	var just [][]string
	var err error

	f := func(s [][]string) maybe.AoAoS { return maybe.JustAoAoS(s[1:]) }

	// Bind AoAoS to AoAoS; good path
	got = good.Bind(f)
	just, err = got.Unbox()
	is.Equal(just, [][]string{[]string{"Goodbye", "Cruel World"}})
	is.Nil(err)

	// Bind AoAoS to AoAoS; bad path
	got = bad.Bind(f)
	is.True(got.IsErr())
}

func TestAoAoSJoin(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]string{
		[]string{"Hello", "World"},
		[]string{"Goodbye", "Cruel World"},
	}
	good, bad := getAoAoSFixtures(input)
	var got maybe.AoS
	var err error

	f := func(x []string) maybe.S { return maybe.JustS(strings.Join(x, " ")) }

	// Join AoAoS to AoS; good path
	got = good.Join(f)
	s, err := got.Unbox()
	is.Equal(s, []string{"Hello World", "Goodbye Cruel World"})
	is.Nil(err)

	// Join AoAoS to AoS; bad path
	got = bad.Join(f)
	is.True(got.IsErr())

	// Join where input is invalid
	badJoin := good.Join(func(x []string) maybe.S { return maybe.ErrS(errors.New("bad int")) })
	is.True(badJoin.IsErr())
}

func TestAoAoSMap(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]string{
		[]string{"Hello", "World"},
		[]string{"Goodbye", "Cruel World"},
	}
	good, bad := getAoAoSFixtures(input)
	var just [][]string
	var err error

	// Map where everything succeeds
	// Function returns first element only
	f := func(xs []string) maybe.AoS { return maybe.JustAoS(xs[0:1]) }
	firsts := good.Map(f)
	just, err = firsts.Unbox()
	is.Equal(just, [][]string{[]string{"Hello"}, []string{"Goodbye"}})
	is.Nil(err)

	// Map where input is invalid
	badInput := bad.Map(f)
	is.True(badInput.IsErr())

	// Map where function returns invalid
	badMap := good.Map(func(xs []string) maybe.AoS { return maybe.ErrAoS(errors.New("bad string")) })
	is.True(badMap.IsErr())
}

func TestAoAoSFlatten(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]string{
		[]string{"Hello", "World"},
		[]string{"Goodbye", "Cruel World"},
	}
	good, bad := getAoAoSFixtures(input)
	var got maybe.AoS
	var just []string
	var err error

	// Good path
	got = good.Flatten()
	just, err = got.Unbox()
	is.Equal(just, []string{"Hello", "World", "Goodbye", "Cruel World"})
	is.Nil(err)

	// Bad path
	got = bad.Flatten()
	is.True(got.IsErr())
}

func TestAoAoSToInt(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	input := [][]string{
		[]string{"42", "23"},
		[]string{"7", "13"},
	}
	good, bad := getAoAoSFixtures(input)
	notNum := maybe.JustAoAoS([][]string{
		[]string{"23", "forty-two"},
		[]string{"nine", "9"},
	})
	var got maybe.AoAoI
	var err error

	f := func(s string) maybe.I { return maybe.NewI(strconv.Atoi(s)) }

	// Convert S to I; good path
	got = good.ToInt(f)
	x, err := got.Unbox()
	is.Equal(x, [][]int{[]int{42, 23}, []int{7, 13}})
	is.Nil(err)

	// Convert S to I; bad path
	got = notNum.ToInt(f)
	is.True(got.IsErr())

	// Convert invalid S to I
	got = bad.ToInt(f)
	is.True(got.IsErr())
}
