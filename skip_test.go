package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	ta := require.New(t)

	s := New()
	ta.Equal(1, len(s.head.pointers))
	ta.Equal(0, len(s.tail.pointers))
	ta.Equal(&s.tail, s.head.pointers[0])
	ta.Equal("\x00", string(s.head.kv))
	ta.Equal("\x02", string(s.tail.kv))
}

func Test_randLvl(t *testing.T) {

	ta := require.New(t)
	ratio := 0.5

	n := 1000 * 1000
	sample := make([]int, 64)
	for i := 0; i < n; i++ {
		v := randLvl(ratio)
		sample[v]++
	}

	fmt.Println(sample[:4])
	for i := 0; i < 4; i++ {
		got := float64(sample[i+1]) / float64(sample[i])
		ta.InDelta(ratio, got, float64(ratio)/5, "%d %d", sample[i], sample[i+1])
	}
}

func TestNewElt_and_Key_Value(t *testing.T) {
	ta := require.New(t)

	e := NewElt(3, "foo", "bar")
	ta.Equal([]*Elt{nil, nil, nil, nil}, e.pointers)
	ta.Equal("\x01foobar", string(e.kv))
	ta.Equal(3, e.keyLen)

	ta.Equal("foo", e.Key())
	ta.Equal("bar", e.Value())
}

func TestElt_Less(t *testing.T) {
	ta := require.New(t)

	a := NewElt(5, "aa", "bb")
	a2 := NewElt(3, "aa", "cc")
	b := NewElt(3, "ab", "cc")

	ta.False(a.Less(a))

	ta.True(a.Less(b))
	ta.False(a.Less(a2))
	ta.False(a2.Less(a))
	ta.False(b.Less(a))

	s := New()
	ta.True(s.head.Less(a))
	ta.True(a.Less(&s.tail))
}

func TestElt_String(t *testing.T) {
	ta := require.New(t)
	a := NewElt(5, "aa", "bb")
	ta.Equal("E:aa=bb", a.String())

	s := New()
	ta.Equal("H:=", s.head.String())
	ta.Equal("T:=", s.tail.String())
}

func TestSkiplist_String(t *testing.T) {

	ta := require.New(t)

	s := New()

	a := NewElt(5, "aa", "bb")
	s.head.pointers[0] = a
	a.pointers[0] = &s.tail

	got := s.String()
	ta.Equal("H:= > E:aa=bb > T:=", got)
}

func TestSkiplist_Add_searchElt_Remove_Get(t *testing.T) {

	ta := require.New(t)

	s := New()
	s.Add("a", "")
	ta.Equal("H:= > E:a= > T:=", s.String())

	s.Add("b", "")
	ta.Equal("H:= > E:a= > E:b= > T:=", s.String())

	got := s.Add("c", "")
	ta.False(got)
	ta.Equal("H:= > E:a= > E:b= > E:c= > T:=", s.String())

	got = s.Add("b", "newB")
	ta.True(got)
	ta.Equal("H:= > E:a= > E:b=newB > E:c= > T:=", s.String())

	got = s.Remove("a")
	ta.True(got)
	ta.Equal("H:= > E:b=newB > E:c= > T:=", s.String())

	got = s.Remove("a")
	ta.False(got)
	ta.Equal("H:= > E:b=newB > E:c= > T:=", s.String())

	v, found := s.Get("a")
	ta.Equal("", v)
	ta.False(found)

	v, found = s.Get("b")
	ta.Equal("newB", v)
	ta.True(found)
}
