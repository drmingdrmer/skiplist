package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
)

var (
	// the chance to add next level.
	// a classic impl is 1/2.
	// usually 1/4 performs best
	defaultRatio = float64(0.3)
)

type Node struct {
	nexts []*Node

	b uint32
	// there are 2 seg in  k:
	// the first byte is node type: head: 00, tail: 02, other node: 01
	// so that head always less than other node and other nodes are always less
	// than tail.
	k []byte
	v interface{}
}

type SkipList struct {
	head, tail Node
}

func New() *SkipList {
	s := &SkipList{}
	s.head.nexts = make([]*Node, 1)
	s.head.nexts[0] = &s.tail
	s.head.k = []byte{0}
	s.tail.k = []byte{2}

	return s
}

// func strToU32(k string) uint32 {
//     l := len(k)
//     switch l {
//     case 0:
//         return 0
//     case 1:
//         return uint32(k[0]) << 24
//     case 2:
//         return uint32(k[0])<<24 | uint32(k[1])<<16
//     case 3:
//         return uint32(k[0])<<24 | uint32(k[1])<<16 | uint32(k[2])<<8
//     default:
//         return uint32(k[0])<<24 | uint32(k[1])<<16 | uint32(k[2])<<8 | uint32(k[3])
//     }
// }

func NewNode(lvl int, k string, v interface{}) *Node {

	e := &Node{
		nexts: make([]*Node, lvl+1),
		// b:     strToU32(k),
		k: make([]byte, 1+len(k)),
		v: v,
	}
	// \x01 indicates it is a payload node.
	e.k[0] = 1
	copy(e.k[1:], k)

	return e
}

func (e *Node) Key() string {
	return string(e.k[1:])
}

func (e *Node) Value() interface{} {
	return e.v
}

// Less return true if e < o.
// e.Less(o) implies !o.Less(e)
// !e.Less(o) && ! o.Less(e) implies e == o
func (e *Node) Less(o *Node) bool {
	// if e.k[0] != o.k[0] {
	//     return e.k[0] < o.k[0]
	// }

	// return e.b < o.b || (e.b == o.b &&
	//     bytes.Compare(e.k, o.k) < 0)
	return bytes.Compare(e.k, o.k) < 0
}

func (e *Node) String() string {
	typ := e.k[0]
	if typ == 0 {
		return "H"
	}
	if typ == 2 {
		return "T"
	}
	return e.Key() + "=" + fmt.Sprintf("%s", e.v)
}

func (s *SkipList) String() string {

	var res []string
	cur := &s.head
	for ; cur != &s.tail; cur = cur.nexts[0] {
		res = append(res, cur.String())
	}
	res = append(res, s.tail.String())
	return strings.Join(res, " > ")
}

func (s *SkipList) DebugStr() string {
	var lines []string

	maxLvl := len(s.head.nexts) - 1
	for cur := &s.head; cur != &s.tail; cur = cur.nexts[0] {
		var l []string
		for lvl := maxLvl; lvl >= 0; lvl-- {
			if lvl >= len(cur.nexts) {
				l = append(l, "| ")
				continue
			}

			l = append(l, "v ")
		}
		l = append(l, cur.String())
		lines = append(lines, strings.Join(l, ""))
	}

	return strings.Join(lines, "\n")
}

// Add a new kv-pair into SkipList, returns true if it overrides an existent
// record.
func (s *SkipList) Add(k, v string) bool {
	lvl := randLevel(defaultRatio)
	if lvl > len(s.head.nexts) {
		lvl = len(s.head.nexts)
	}

	e := NewNode(lvl, k, v)
	for len(e.nexts) > len(s.head.nexts) {
		s.head.nexts = append(s.head.nexts, &s.tail)
	}

	ps, equal := s.searchNode(e)

	if equal {
		s.removeNode(ps)
	}

	for i := 0; i < len(e.nexts); i++ {
		e.nexts[i] = ps[i].nexts[i]
		ps[i].nexts[i] = e
	}

	return equal
}

// Remove a new k from SkipList, returns true if removed.
func (s *SkipList) Remove(k string) bool {

	e := NewNode(0, k, "")

	ps, equal := s.searchNode(e)

	if equal {
		s.removeNode(ps)
	}

	return equal
}

type KVer interface {
	Key() string
	Value() interface{}
}

// Find the first node e so that k <= e.
// It returns an interface that KVer and bool indicate if
// k == e.
// KVer is nil if there is no node that k <= e.
func (s *SkipList) Get(k string) (KVer, bool) {

	e := NewNode(0, k, "")
	ps, equal := s.searchNode(e)

	nxt := ps[0].nexts[0]
	if nxt == &s.tail {
		return nil, false
	}

	return nxt, equal
}

func (s *SkipList) removeNode(ps []*Node) {
	p := ps[0].nexts[0]
	for i := 0; i < len(ps); i++ {
		if ps[i].nexts[i] == p {
			ps[i].nexts[i] = p.nexts[i]
		}
	}
}

// Search for k and returns a slice of *Node
// and if the previous node exactly equal e.
//
//   3   ------------------------>
//   2   -------------------> 2 ->
//   1   ---------> 1 ------> 1 ->
//   0 head -> a -> b -> c -> d -> tail
//
func (s *SkipList) searchNode(e *Node) ([]*Node, bool) {

	h := &s.head

	lvl := len(h.nexts) - 1
	rst := make([]*Node, lvl+1)

	cur := h
	for ; lvl >= 0; lvl-- {

		// find the first cur that e <= cur.next
		for cur.nexts[lvl].Less(e) {
			cur = cur.nexts[lvl]
		}

		rst[lvl] = cur
	}

	nxt := cur.nexts[0]

	return rst, !e.Less(nxt)
}

func randLevel(ratio float64) int {
	x := 0
	for rand.Float64() < ratio {
		x++
	}
	return x
}
