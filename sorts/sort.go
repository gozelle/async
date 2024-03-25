package sorts

import (
	"sort"

	"github.com/gozelle/async/generics"
)

var _ sort.Interface = (*sorter[any])(nil)

type sorter[T any] struct {
	list []T
	less func(a, b T) bool
}

func (s sorter[T]) Len() int {
	return len(s.list)
}

func (s sorter[T]) Less(i, j int) bool {
	return s.less(s.list[i], s.list[j])
}

func (s sorter[T]) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func Sort[T any](items []T, less func(a, b T) bool) {
	s := sorter[T]{
		list: items,
		less: less,
	}
	sort.Sort(s)
}

func SortAsc[T any, K generics.Ordered](items []T, key func(item T) K) {
	s := &sorter[T]{
		list: items,
		less: func(a, b T) bool {
			return key(a) < key(b)
		},
	}
	sort.Sort(s)
}

func SortDesc[T any, K generics.Ordered](items []T, key func(item T) K) {
	s := &sorter[T]{
		list: items,
		less: func(a, b T) bool {
			return key(a) > key(b)
		},
	}
	sort.Sort(s)
}
