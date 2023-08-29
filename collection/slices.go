package collection

// UniqueKeys Slices Get a unique slices
// T => slice item type
// K => slice item filter key type
func UniqueKeys[T any, K comparable](items []T, fn func(item T) K) []K {
	r := make([]K, 0)
	m := map[K]struct{}{}
	for _, v := range items {
		k := fn(v)
		if _, ok := m[k]; !ok {
			m[k] = struct{}{}
			r = append(r, k)
		}
	}
	return r
}

// UniqueMap convert slice to a map by unique key
// T => slice item type
// K => slice item filter key type
func UniqueMap[T any, K comparable](items []T, fn func(item T) (K, T)) map[K]T {
	m := map[K]T{}
	for _, v := range items {
		k, vv := fn(v)
		if _, ok := m[k]; !ok {
			m[k] = vv
		}
	}
	return m
}
