package treemap

import (
	"fmt"
	"github.com/pobochiigo/gotil/comparator"
	"github.com/pobochiigo/gotil/core"
	"golang.org/x/exp/constraints"
	"reflect"
)

type treeMap[K constraints.Ordered, V any] struct {
	root       *entry[K, V]
	size       int
	modCount   int
	comparator comparator.Comparable[K]
	entries    []core.Entry[K, V]
	keys       []K
	values     []V
}

func New[K constraints.Ordered, V any](comparator comparator.Comparable[K]) core.Map[K, V] {
	return &treeMap[K, V]{comparator: comparator}
}

func (m *treeMap[K, V]) Size() int {
	return m.size
}

func (m *treeMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

func (m *treeMap[K, V]) Clear() {
	m.modCount++
	m.size = 0
	m.root = nil
}

func (m *treeMap[K, V]) ContainsKey(key K) bool {
	return m.getEntry(key) != nil
}

func (m *treeMap[K, V]) ContainsValue(value V) bool {
	for e := m.getFirstEntry(); e != nil; e = m.successor(e) {
		if reflect.DeepEqual(value, e.value) {
			return true
		}
	}
	return false
}

func (m *treeMap[K, V]) Get(key K) *V {
	if e := m.getEntry(key); e != nil {
		return &e.value
	}
	return nil
}

func (m *treeMap[K, V]) Put(key K, value V) *V {
	t := m.root
	if t == nil {
		m.compare(key, key) // type (and possibly null) check
		m.root = m.newEntry(key, value, nil)
		m.size = 1
		m.modCount++
		return nil
	}
	var cmp int
	var parent *entry[K, V]
	// split comparator and comparable paths
	cpr := m.comparator
	if cpr != nil {
		for t != nil {
			parent = t
			cmp = cpr(key, t.key)
			switch {
			case cmp < 0:
				t = t.left
			case cmp > 0:
				t = t.right
			default:
				return t.SetValue(value)
			}
		}
	} else {
		//if key == nil {
		//	return nil, errors.New("nil pointer exception")
		//}
		for t != nil {
			parent = t
			cmp = compareKeys(key, t.key)
			switch {
			case cmp < 0:
				t = t.left
			case cmp > 0:
				t = t.right
			default:
				return t.SetValue(value)
			}
		}
	}
	e := m.newEntry(key, value, parent)
	if cmp < 0 {
		parent.left = e
	} else {
		parent.right = e
	}
	m.fixAfterInsertion(e)
	m.size++
	m.modCount++
	return nil
}

type Key[K constraints.Ordered] struct {
	value K
}

func (k *Key[K]) compareTo(key Key[K]) int {
	switch {
	case k.value < key.value:
		return -1
	case k.value > key.value:
		return 1
	default:
		return 0
	}
}

func compareKeys[K constraints.Ordered](k1, k2 K) int {
	switch {
	case k1 < k2:
		return -1
	case k1 > k2:
		return 1
	default:
		return 0
	}
}

func (m *treeMap[K, V]) Remove(key K) *V {
	if e := m.getEntry(key); e != nil {
		oldValue := e.value
		m.deleteEntry(e)
		return &oldValue
	}
	return nil
}

// PutAll todo: consider move to AbstractMap
func (m *treeMap[K, V]) PutAll(p core.Map[K, V]) {
	for _, e := range p.Entries() {
		m.Put(e.Key(), e.Value())
	}
}

func (m *treeMap[K, V]) PutMAll(p map[K]V) {
	for k, v := range p {
		m.Put(k, v)
	}
}

func (m *treeMap[K, V]) Keys() []K {
	if ks := m.keys; ks != nil {
		return ks
	}
	m.newKeys()
	return m.keys
}

func (m *treeMap[K, V]) Values() []V {
	if vs := m.values; vs != nil {
		return vs
	}
	m.newValues()
	return m.values
}

func (m *treeMap[K, V]) Entries() []core.Entry[K, V] {
	if es := m.entries; es != nil {
		return es
	}
	m.newEntries()
	return m.entries
}

func (m *treeMap[K, V]) ToMap() map[K]V {
	mp := make(map[K]V, m.size)
	for _, e := range m.Entries() {
		mp[e.Key()] = e.Value()
	}
	return mp
}

func (m *treeMap[K, V]) newEntry(key K, value V, parent *entry[K, V]) *entry[K, V] {
	return &entry[K, V]{key: key, value: value, parent: parent, color: BLACK}
}

// Compares two keys using the correct comparison method for this treeMap.
func (m *treeMap[K, V]) compare(k1, k2 K) int {
	if m.comparator == nil {
		return compareKeys(k1, k2)
	}
	return m.comparator(k1, k2)
}

func (m *treeMap[K, V]) fixAfterInsertion(e *entry[K, V]) {
	e.color = RED
	for e != nil && e != m.root && e.parent.color == RED {
		if m.parentOf(e) == m.leftOf(m.parentOf(m.parentOf(e))) {
			y := m.rightOf(m.parentOf(m.parentOf(e)))
			if m.colorOf(y) == RED {
				m.setColor(m.parentOf(e), BLACK)
				m.setColor(y, BLACK)
				m.setColor(m.parentOf(m.parentOf(e)), RED)
				e = m.parentOf(m.parentOf(e))
			} else {
				if e == m.rightOf(m.parentOf(e)) {
					e = m.parentOf(e)
					m.rotateLeft(e)
				}
				m.setColor(m.parentOf(e), BLACK)
				m.setColor(m.parentOf(m.parentOf(e)), RED)
				m.rotateRight(m.parentOf(m.parentOf(e)))
			}
		} else {
			y := m.leftOf(m.parentOf(m.parentOf(e)))
			if m.colorOf(y) == RED {
				m.setColor(m.parentOf(e), BLACK)
				m.setColor(y, BLACK)
				m.setColor(m.parentOf(m.parentOf(e)), RED)
				e = m.parentOf(m.parentOf(e))
			} else {
				if e == m.leftOf(m.parentOf(e)) {
					e = m.parentOf(e)
					m.rotateRight(e)
				}
				m.setColor(m.parentOf(e), BLACK)
				m.setColor(m.parentOf(m.parentOf(e)), RED)
				m.rotateLeft(m.parentOf(m.parentOf(e)))
			}
		}
	}
	m.root.color = BLACK
}

func (m *treeMap[K, V]) parentOf(e *entry[K, V]) *entry[K, V] {
	if e == nil {
		return nil
	}
	return e.parent
}

func (m *treeMap[K, V]) leftOf(e *entry[K, V]) *entry[K, V] {
	if e == nil {
		return nil
	}
	return e.left
}

func (m *treeMap[K, V]) rightOf(e *entry[K, V]) *entry[K, V] {
	if e == nil {
		return nil
	}
	return e.right
}

func (m *treeMap[K, V]) colorOf(e *entry[K, V]) color {
	if e == nil {
		return BLACK
	}
	return e.color
}

func (m *treeMap[K, V]) setColor(e *entry[K, V], c color) {
	if e != nil {
		e.color = c
	}
}

func (m *treeMap[K, V]) rotateLeft(e *entry[K, V]) {
	if e != nil {
		r := e.right
		e.right = r.left
		if r.left != nil {
			r.left.parent = e
		}
		r.parent = e.parent
		if e.parent == nil {
			m.root = r
		} else if e.parent.left == e {
			e.parent.left = r
		} else {
			e.parent.right = r
		}
		r.left = e
		e.parent = r
	}
}

func (m *treeMap[K, V]) rotateRight(e *entry[K, V]) {
	if e != nil {
		l := e.left
		e.left = l.right
		if l.right != nil {
			l.right.parent = e
		}
		l.parent = e.parent
		if e.parent == nil {
			m.root = l
		} else if e.parent.right == e {
			e.parent.right = l
		} else {
			e.parent.left = l
		}
		l.right = e
		e.parent = l
	}
}

func (m *treeMap[K, V]) getEntry(key K) *entry[K, V] {
	// Offload comparator-based version for sake of performance
	if m.comparator != nil {
		return m.getEntryUsingComparator(key)
	}
	p := m.root
	for p != nil {
		cmp := compareKeys(key, p.key)
		if cmp < 0 {
			p = p.left
		} else if cmp > 0 {
			p = p.right
		} else {
			return p
		}
	}
	return nil
}

func (m *treeMap[K, V]) getEntryUsingComparator(key K) *entry[K, V] {
	cpr := m.comparator
	if cpr != nil {
		p := m.root
		for p != nil {
			cmp := cpr(key, p.key)
			if cmp < 0 {
				p = p.left
			} else if cmp > 0 {
				p = p.right
			} else {
				return p
			}
		}
	}
	return nil
}

func (m *treeMap[K, V]) getFirstEntry() *entry[K, V] {
	p := m.root
	if p != nil {
		for p.left != nil {
			p = p.left
		}
	}
	return p
}

func (m *treeMap[K, V]) getLastEntry() *entry[K, V] {
	p := m.root
	if p != nil {
		for p.right != nil {
			p = p.right
		}
	}
	return p
}

func (m *treeMap[K, V]) successor(e *entry[K, V]) *entry[K, V] {
	if e == nil {
		return nil
	} else if e.right != nil {
		p := e.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		p := e.parent
		ch := e
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}
}

func (m *treeMap[K, V]) predecessor(e *entry[K, V]) *entry[K, V] {
	if e == nil {
		return nil
	} else if e.left != nil {
		p := e.left
		for p.right != nil {
			p = p.right
		}
		return p
	} else {
		p := e.parent
		ch := e
		for p != nil && ch == p.left {
			ch = p
			p = p.parent
		}
		return p
	}
}

func (m *treeMap[K, V]) deleteEntry(e *entry[K, V]) {
	m.modCount++
	m.size--

	// If strictly internal, copy successor's element to e and then make e
	// point to successor.
	if e.left != nil && e.right != nil {
		s := m.successor(e)
		e.key = s.key
		e.value = s.value
		e = s
	} // e has 2 children

	// Start fixup at replacement node, if it exists.
	var repl *entry[K, V]
	if e.left != nil {
		repl = e.left
	} else {
		repl = e.right
	}
	if repl != nil {
		// Link replacement to parent
		repl.parent = e.parent
		if e.parent == nil {
			m.root = repl
		} else if e == e.parent.left {
			e.parent.left = repl
		} else {
			e.parent.right = repl
		}

		// Nil out links, so they are OK to use by fixAfterDeletion.
		e.left = nil
		e.right = nil
		e.parent = nil

		// Fix replacement
		if e.color == BLACK {
			m.fixAfterDeletion(repl)
		}
	} else if e.parent == nil { // return if we are the only node.
		m.root = nil
	} else { //  No children. Use self as phantom replacement and unlink.
		if e.color == BLACK {
			m.fixAfterDeletion(e)
		}

		if e.parent != nil {
			if e == e.parent.left {
				e.parent.left = nil
			} else if e == e.parent.right {
				e.parent.right = nil
			}
			e.parent = nil
		}
	}
}

func (m *treeMap[K, V]) fixAfterDeletion(e *entry[K, V]) {
	for e != m.root && m.colorOf(e) == BLACK {
		if e == m.leftOf(m.parentOf(e)) {
			sib := m.rightOf(m.parentOf(e))

			if m.colorOf(sib) == RED {
				m.setColor(sib, BLACK)
				m.setColor(m.parentOf(e), RED)
				m.rotateLeft(m.parentOf(e))
				sib = m.rightOf(m.parentOf(e))
			}

			if m.colorOf(m.leftOf(sib)) == BLACK &&
				m.colorOf(m.rightOf(sib)) == BLACK {
				m.setColor(sib, RED)
				e = m.parentOf(e)
			} else {
				if m.colorOf(m.rightOf(sib)) == BLACK {
					m.setColor(m.leftOf(sib), BLACK)
					m.setColor(sib, RED)
					m.rotateRight(sib)
					sib = m.rightOf(m.parentOf(e))
				}
				m.setColor(sib, m.colorOf(m.parentOf(e)))
				m.setColor(m.parentOf(e), BLACK)
				m.setColor(m.rightOf(sib), BLACK)
				m.rotateLeft(m.parentOf(e))
				e = m.root
			}
		} else { // symmetric
			sib := m.leftOf(m.parentOf(e))

			if m.colorOf(sib) == RED {
				m.setColor(sib, BLACK)
				m.setColor(m.parentOf(e), RED)
				m.rotateRight(m.parentOf(e))
				sib = m.leftOf(m.parentOf(e))
			}

			if m.colorOf(m.rightOf(sib)) == BLACK &&
				m.colorOf(m.leftOf(sib)) == BLACK {
				m.setColor(sib, RED)
				e = m.parentOf(e)
			} else {
				if m.colorOf(m.leftOf(sib)) == BLACK {
					m.setColor(m.rightOf(sib), BLACK)
					m.setColor(sib, RED)
					m.rotateLeft(sib)
					sib = m.leftOf(m.parentOf(e))
				}
				m.setColor(sib, m.colorOf(m.parentOf(e)))
				m.setColor(m.parentOf(e), BLACK)
				m.setColor(m.leftOf(sib), BLACK)
				m.rotateRight(m.parentOf(e))
				e = m.root
			}
		}
	}

	m.setColor(e, BLACK)
}

func (m *treeMap[K, V]) newEntries() {
	first := m.getFirstEntry()
	m.entries = make([]core.Entry[K, V], 0, m.size)
	m.entries = append(m.entries, first)
	s := m.successor(first)
	for s != nil {
		m.entries = append(m.entries, s)
		s = m.successor(s)
	}
}

func (m *treeMap[K, V]) newKeys() {
	m.keys = make([]K, 0, m.size)
	for _, e := range m.Entries() {
		m.keys = append(m.keys, e.Key())
	}
}

func (m *treeMap[K, V]) newValues() {
	m.values = make([]V, 0, m.size)
	for _, e := range m.Entries() {
		m.values = append(m.values, e.Value())
	}
}

type color bool

const (
	BLACK color = true
	RED   color = false
)

type entry[K comparable, V any] struct {
	key    K
	value  V
	left   *entry[K, V]
	right  *entry[K, V]
	parent *entry[K, V]
	color  color // true - black, false - red
}

func (e *entry[K, V]) Key() K {
	return e.key
}

func (e *entry[K, V]) Value() V {
	return e.value
}

func (e *entry[K, V]) SetValue(value V) *V {
	oldValue := e.value
	e.value = value
	return &oldValue
}

func (e *entry[K, V]) String() string {
	return fmt.Sprintf("%v=%v", e.key, e.value)
}
