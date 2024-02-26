package core

import "golang.org/x/exp/constraints"

// Map interface definition
type Map[K constraints.Ordered, V any] interface {
	Size() int
	IsEmpty() bool
	ContainsKey(key K) bool
	ContainsValue(value V) bool
	Get(key K) *V
	Put(key K, value V) *V
	Remove(key K) *V
	PutAll(p Map[K, V])
	PutMAll(p map[K]V)
	Clear()
	Keys() []K
	Values() []V
	Entries() []Entry[K, V]
	ToMap() map[K]V
}

type Entry[K constraints.Ordered, V any] interface {
	Key() K
	Value() V
	SetValue(value V) *V
}
