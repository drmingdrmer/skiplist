package skiplist

import (
	"math/rand"
	"strings"
)

var (
	// the chance to add next level.
	// a classic impl is 1/2.
	// usually 1/4 perform best
	ratio = float64(0.5)
)

type Elt struct {
	pointers []*Elt

	keyLen int
	// there are 3 seg in  kv:
	// the first byte is elt type: head: 00, tail: 02, other elt: 01
	// so that head always less than other node and other elts are always less
	// than tail.
	kv []byte
}

type Skiplist struct {
	head, tail Elt
}

func New() *Skiplist {
	s := &Skiplist{}
	s.head.pointers = make([]*Elt, 1)
	s.head.pointers[0] = &s.tail
	s.head.kv = []byte{0}
	s.tail.kv = []byte{2}

	return s
}

func NewElt(lvl int, k, v string) *Elt {

	e := &Elt{
		pointers: make([]*Elt, lvl+1),
		keyLen:   len(k),
		// \x01 indicates it is a payload elt.
		kv: []byte("\x01" + k + v),
	}

	return e
}

func (e *Elt) Key() string {
	return string(e.kv[1 : 1+e.keyLen])
}

func (e *Elt) cmpKey() string {
	return string(e.kv[:1+e.keyLen])
}

func (e *Elt) Value() string {
	return string(e.kv[1+e.keyLen:])
}

// Less return true if e < o.
// e.Less(o) implies !o.Less(e)
// !e.Less(o) && ! o.Less(e) implies e == o
func (e *Elt) Less(o *Elt) bool {
	ka := e.cmpKey()
	kb := o.cmpKey()
	return ka < kb
}

func (e *Elt) String() string {
	typs := map[byte]string{
		0: "H:", // head
		1: "E:", // normal elt
		2: "T:", // tail
	}
	typ := e.kv[0]
	return typs[typ] + e.Key() + "=" + e.Value()
}

func (s *Skiplist) String() string {

	res := []string{}
	cur := &s.head
	for ; cur != &s.tail; cur = cur.pointers[0] {
		res = append(res, cur.String())

	}
	res = append(res, s.tail.String())
	return strings.Join(res, " > ")
}

// Add a new kv-pair into skiplist, returns true if it overrides an existent
// record.
func (s *Skiplist) Add(k, v string) bool {
	lvl := randLvl(ratio)

	e := NewElt(lvl, k, v)
	for len(e.pointers) > len(s.head.pointers) {
		s.head.pointers = append(s.head.pointers, &s.tail)
	}

	ps, equal := s.searchElt(e)

	if equal {
		s.removeElt(ps)
	}

	for i := 0; i < len(e.pointers); i++ {
		e.pointers[i] = ps[i].pointers[i]
		ps[i].pointers[i] = e
	}

	return equal
}

// Remove a new k from skiplist, returns true if removed.
func (s *Skiplist) Remove(k string) bool {

	e := NewElt(0, k, "")

	ps, equal := s.searchElt(e)

	if equal {
		s.removeElt(ps)
	}

	return equal
}

func (s *Skiplist) Get(k string) (string, bool) {

	e := NewElt(0, k, "")
	ps, equal := s.searchElt(e)

	if equal {
		return ps[0].pointers[0].Value(), true
	}
	return "", false
}

func (s *Skiplist) removeElt(ps []*Elt) {
	p := ps[0].pointers[0]
	for i := 0; i < len(ps); i++ {
		if ps[i].pointers[i] == p {
			ps[i].pointers[i] = p.pointers[i]
		}
	}
}

// Search for k and returns a slice of *Elt
// and if the previous elt exactly equal e.
//
//   3   ------------------------>
//   2   -------------------> 2 ->
//   1   ---------> 1 ------> 1 ->
//   0 head -> a -> b -> c -> d -> tail
//
func (s *Skiplist) searchElt(e *Elt) ([]*Elt, bool) {

	h := &s.head

	lvl := len(h.pointers) - 1
	rst := make([]*Elt, lvl+1)

	cur := h
	for ; lvl >= 0; lvl-- {

		// find the first cur that e <= cur.next
		for cur.pointers[lvl].Less(e) {
			cur = cur.pointers[lvl]
		}

		rst[lvl] = cur
	}

	nxt := cur.pointers[0]

	return rst, !nxt.Less(e) && !e.Less(nxt)
}

func randLvl(ratio float64) int {
	x := 0
	for rand.Float64() < ratio {
		x++
	}
	return x
}
