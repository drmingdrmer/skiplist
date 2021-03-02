package skiplist

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	ta := require.New(t)

	s := New()
	ta.Equal(1, len(s.head.nexts))
	ta.Equal(0, len(s.tail.nexts))
	ta.Equal(&s.tail, s.head.nexts[0])
	ta.Equal("\x00", string(s.head.k))
	ta.Equal("\x02", string(s.tail.k))
}

func Test_randLvl(t *testing.T) {

	ta := require.New(t)
	ratio := 0.5

	n := 1000 * 1000
	sample := make([]int, 64)
	for i := 0; i < n; i++ {
		v := randLevel(ratio)
		sample[v]++
	}

	// fmt.Println(sample[:4])
	for i := 0; i < 4; i++ {
		got := float64(sample[i+1]) / float64(sample[i])
		ta.InDelta(ratio, got, float64(ratio)/5, "%d %d", sample[i], sample[i+1])
	}
}

func TestNewElt_and_Key_Value(t *testing.T) {
	ta := require.New(t)

	e := NewNode(3, "foo", "bar")
	ta.Equal([]*Node{nil, nil, nil, nil}, e.nexts)
	ta.Equal("\x01foo", string(e.k))

	ta.Equal("foo", e.Key())
	ta.Equal("bar", e.Value())
}

func TestElt_Less(t *testing.T) {
	ta := require.New(t)

	a := NewNode(5, "aa", "bb")
	a2 := NewNode(3, "aa", "cc")
	b := NewNode(3, "ab", "cc")

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
	a := NewNode(5, "aa", "bb")
	ta.Equal("aa=bb", a.String())

	s := New()
	ta.Equal("H", s.head.String())
	ta.Equal("T", s.tail.String())
}

func TestSkiplist_String(t *testing.T) {

	ta := require.New(t)

	s := New()

	a := NewNode(5, "aa", "bb")
	s.head.nexts[0] = a
	a.nexts[0] = &s.tail

	got := s.String()
	ta.Equal("H > aa=bb > T", got)
}

func TestSkiplist_Add_searchElt_Remove_Get(t *testing.T) {

	ta := require.New(t)

	s := New()
	s.Add("a", "")
	ta.Equal("H > a= > T", s.String())

	s.Add("b", "")
	ta.Equal("H > a= > b= > T", s.String())

	got := s.Add("c", "")
	ta.False(got)
	ta.Equal("H > a= > b= > c= > T", s.String())

	got = s.Add("c", "newC")
	ta.True(got)
	ta.Equal("H > a= > b= > c=newC > T", s.String())

	got = s.Add("b", "newB")
	ta.True(got)
	ta.Equal("H > a= > b=newB > c=newC > T", s.String())

	got = s.Remove("a")
	ta.True(got)
	ta.Equal("H > b=newB > c=newC > T", s.String())

	got = s.Remove("a")
	ta.False(got)
	ta.Equal("H > b=newB > c=newC > T", s.String())

	kv, found := s.Get("a")
	ta.Equal("b", kv.Key())
	ta.False(found)

	kv, found = s.Get("b")
	ta.Equal("newB", kv.Value())
	ta.True(found)

	kv, found = s.Get("d")
	ta.Nil(kv)
	ta.False(found)
}

func TestSkiplist_DebugStr(t *testing.T) {
	ta := require.New(t)
	s := New()
	letters := []byte("abcdefghijkl")

	n := 15
	for i := 0; i < n; i++ {
		x := rand.Int() % len(letters)
		y := rand.Int() % len(letters)
		k := string([]byte{letters[x], letters[y]})
		s.Add(k, k)
	}

	// fmt.Println(s.DebugStr())

	want := `
v v v H
| | v ab=ab
| | v ae=ae
| | v ba=ba
| | v bh=bh
| | v cc=cc
| v v dc=dc
v v v ei=ei
| | v ge=ge
| | v gh=gh
| | v hc=hc
| | v hi=hi
| | v ic=ic
| | v kh=kh
| | v kj=kj`
	ta.Equal(want[1:], s.DebugStr())
}
