package core

import (
	"golang.org/x/exp/constraints"
)

type AbstractMap[K constraints.Ordered, V any] struct {
}

func (a *AbstractMap[K, V]) Size() int {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) ContainsKey(key K) bool {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) ContainsValue(value V) bool {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Get(key K) V {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Put(key K, value V) V {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Remove(value V) V {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) PutAll(p Map[K, V]) {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Clear() {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Keys() []K {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Values() []V {
	//TODO implement me
	panic("implement me")
}

func (a *AbstractMap[K, V]) Entries() []Entry[K, V] {
	//TODO implement me
	panic("implement me")
}
