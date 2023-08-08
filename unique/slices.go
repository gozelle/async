package unique

// SliceKeys Slices Get a unique slices
// T => slice item type
// K => slice item filter key type
// V => the final result type 
func SliceKeys[T any, K comparable](items []T, kv func(item T) K) []K {
	r := make([]K, 0)
	m := map[K]struct{}{}
	for _, v := range items {
		k := kv(v)
		if _, ok := m[k]; !ok {
			m[k] = struct{}{}
			r = append(r, k)
		}
	}
	return r
}
