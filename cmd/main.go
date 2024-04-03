package main

import (
	"fmt"
	"github.com/pobochiigo/gotil/comparator"
	"github.com/pobochiigo/gotil/treemap"
	"math/rand"
	"time"
)

type item struct {
	value int
}

type tMap map[int]*item

func (tm *tMap) ContainsKey(key int) bool {
	return map[int]*item(*tm)[key] != nil
}

func main() {
	testTreeMap()
	//testMap()
}

const Chars = "ABCDEFGHIJKLMOPQRSTUVWXYZ"

func testTreeMap() {

	cmpr := comparator.New[string](comparator.ASC)
	tm := treemap.NewWithComparator[string, item](cmpr)
	charsSize := len(Chars)

	s := time.Now()

	for i := 0; i < 13; i++ {
		key := string(Chars[rand.Intn(charsSize)]) + string(Chars[rand.Intn(charsSize)])
		tm.Put(key, item{i})
	}

	println("populated by:", time.Now().Sub(s).String())

	entries := tm.Entries()
	fmt.Printf("entries: %v\n", entries)
}

func testMap() {

	m := make(tMap)

	s := time.Now()

	for i := 0; i < 1000000; i++ {
		m[i] = &item{i}
	}

	println("populated by:", time.Now().Sub(s).String())

	s = time.Now()
	println(m.ContainsKey(0))
	println("searched key=0 by:", time.Now().Sub(s).String())
	s = time.Now()
	println(m.ContainsKey(500000))
	println("searched key=500000 by:", time.Now().Sub(s).String())
	s = time.Now()
	println(m.ContainsKey(999999))
	println("searched key=999999 by:", time.Now().Sub(s).String())
	s = time.Now()
	println(m.ContainsKey(1000000))
	println("searched key=1000000 by:", time.Now().Sub(s).String())
}
