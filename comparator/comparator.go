package comparator

import "golang.org/x/exp/constraints"

type Order int

const (
	ASC  Order = 1
	DESC Order = -1
)

type Comparator[K constraints.Ordered] interface {
	Compare(k1, k2 K) int
	Asc()
	Desc()
}

type compare[K constraints.Ordered] struct {
	order Order
}

func (c *compare[K]) Asc() {
	c.order = ASC
}

func (c *compare[K]) Desc() {
	c.order = DESC
}

func (c *compare[K]) Compare(k1, k2 K) int {
	if k1 < k2 {
		return int(-1 * c.order)
	}
	if k1 > k2 {
		return int(1 * c.order)
	}
	return 0
}

func New2[K constraints.Ordered]() Comparator[K] {
	return &compare[K]{ASC}
}
