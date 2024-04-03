package treemap

import (
	"fmt"
	"github.com/pobochiigo/gotil/comparator"
	"github.com/pobochiigo/gotil/core"
	"github.com/emirpasic/gods/maps/treemap"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

const Chars = "abcdefghijklmopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXYZ0123456789"

var charsSize = len(Chars)

type item struct {
	id      string
	value   int
	created time.Time
	valid   bool
	params  map[string]any
}

var tm core.Map[string, item]
var tmc core.Map[string, item]
var gods *treemap.Map
var lock sync.RWMutex

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
	}
}

func BenchmarkMapPut_NoComp(b *testing.B) {
	b.Cleanup(func() {
		b.Logf("map.size=%d", tm.Size())
		tm.Clear()
	})

	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("%s%d", string(Chars[rand.Intn(charsSize)]), i)
		val := item{
			id:      id,
			value:   i,
			created: time.Now(),
			valid:   true,
			params:  map[string]any{"id": id, "value": i, "created": time.Now(), "valid": true},
		}
		tm.Put(id, val)
	}
}

func BenchmarkMapPut_TreeMap(b *testing.B) {
	//b.Cleanup(func() {
	//	b.Logf("map.size=%d", tmc.Size())
	//	tmc.Clear()
	//})
	//
	//for i := 0; i < b.N; i++ {
	//	id := fmt.Sprintf("%s%d", string(Chars[rand.Intn(charsSize)]), i)
	//	val := item{
	//		id:      id,
	//		value:   i,
	//		created: time.Now(),
	//		valid:   true,
	//		params:  map[string]any{"id": id, "value": i, "created": time.Now(), "valid": true},
	//	}
	//	tmc.Put(id, val)
	//}

	b.Cleanup(func() {
		b.Logf("map.size=%d", gods.Size())
		gods.Clear()
	})

	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("%s%d", string(Chars[rand.Intn(charsSize)]), i)
		val := item{
			id:      id,
			value:   i,
			created: time.Now(),
			valid:   true,
			params:  map[string]any{"id": id, "value": i, "created": time.Now(), "valid": true},
		}
		gods.Put(id, val)
	}
}

func BenchmarkMapPut_TreeMap_Sync(b *testing.B) {
	//b.Cleanup(func() {
	//	lock.Lock()
	//	b.Logf("map.size=%d", tmc.Size())
	//	tmc.Clear()
	//	lock.Unlock()
	//})
	//
	//for i := 0; i < b.N; i++ {
	//	go func(i int) {
	//		id := fmt.Sprintf("%s%d", string(Chars[rand.Intn(charsSize)]), i)
	//		val := item{
	//			id:      id,
	//			value:   i,
	//			created: time.Now(),
	//			valid:   true,
	//			params:  map[string]any{"id": id, "value": i, "created": time.Now(), "valid": true},
	//		}
	//
	//		lock.Lock()
	//		tmc.Put(id, val)
	//		lock.Unlock()
	//	}(i)
	//}

	b.Cleanup(func() {
		b.Logf("map.size=%d", gods.Size())
		gods.Clear()
	})

	for i := 0; i < b.N; i++ {
		go func(i int) {
			id := fmt.Sprintf("%s%d", string(Chars[rand.Intn(charsSize)]), i)
			val := item{
				id:      id,
				value:   i,
				created: time.Now(),
				valid:   true,
				params:  map[string]any{"id": id, "value": i, "created": time.Now(), "valid": true},
			}

			lock.Lock()
			gods.Put(id, val)
			lock.Unlock()
		}(i)
	}
}

func TestMain(m *testing.M) {
	println("TestMain::start")
	beforeTest()
	code := m.Run()
	afterTest()
	println("TestMain::finish")
	os.Exit(code)
}

func afterTest() {
	//tm.Clear()
	//tmc.Clear()
	//tm = nil
	//tmc = nil
}

func beforeTest() {
	tm = New[string, item](nil)
	tmc = New[string, item](comparator.New[string](comparator.ASC))
	gods = treemap.NewWithStringComparator()
}
