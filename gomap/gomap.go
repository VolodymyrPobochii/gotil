package gomap

import (
	"github.com/pobochiigo/gotil/comparator"
	"github.com/pobochiigo/gotil/core"
	"golang.org/x/exp/constraints"
	"reflect"
)

type gmap[K constraints.Ordered, V any] struct {
	m map[K]V
	c comparator.Comparable[K]
}

func New[K constraints.Ordered, V any](comparator comparator.Comparable[K]) core.Map[K, V] {
	return &gmap[K, V]{m: make(map[K]V), c: comparator}
}

func (g *gmap[K, V]) Size() int {
	return len(g.m)
}

func (g *gmap[K, V]) IsEmpty() bool {
	return len(g.m) == 0
}

func (g *gmap[K, V]) ContainsKey(key K) bool {
	_, ok := g.m[key]
	return ok
}

func (g *gmap[K, V]) ContainsValue(value V) bool {
	refv := reflect.ValueOf(value)
	for _, v := range g.m {
		if reflect.ValueOf(v) == refv {
			return true
		}
	}
	return false
}

func (g *gmap[K, V]) Get(key K) *V {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) Put(key K, value V) *V {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) Remove(key K) *V {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) PutAll(p core.Map[K, V]) {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) PutMAll(p map[K]V) {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) Clear() {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) Keys() []K {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) Values() []V {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) Entries() []core.Entry[K, V] {
	//TODO implement me
	panic("implement me")
}

func (g *gmap[K, V]) ToMap() map[K]V {
	//TODO implement me
	panic("implement me")
}
