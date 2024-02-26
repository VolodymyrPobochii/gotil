package comparator

import (
	"golang.org/x/exp/constraints"
)

type Comparable[K constraints.Ordered] func(k1, k2 K) int

func (c Comparable[K]) compare(k1, k2 K) int {
	return c(k1, k2)
}

func New[K constraints.Ordered](order Order) Comparable[K] {
	return func(k1, k2 K) int {
		if k1 < k2 {
			return int(-1 * order)
		}
		if k1 > k2 {
			return int(1 * order)
		}
		return 0
	}
}
