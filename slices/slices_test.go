package slices

import "testing"

func TestUniqueKeys(t *testing.T) {
	a := []string{"a", "b", "c", "d", "d", "c", "a"}
	b := UniqueKeys[string, string](a, func(item string) (key string) {
		key = item
		return
	})
	t.Log(b)
	
	type I struct {
		Group string
		User  string
	}
	
	items := []I{
		{Group: "a", User: "tom"},
		{Group: "a", User: "bob"},
		{Group: "b", User: "jack"},
	}
	
	ii := UniqueKeys[I, string](items, func(item I) (key string) {
		return item.Group
	})
	
	t.Log(ii)
}

func TestUniqueMap(t *testing.T) {
	a := []string{"a", "b", "c", "d", "d", "c", "a"}
	b := UniqueKeys[string, string](a, func(item string) (key string) {
		key = item
		return
	})
	t.Log(b)
	
	type I struct {
		Group string
		User  string
	}
	
	items := []I{
		{Group: "a", User: "tom"},
		{Group: "a", User: "bob"},
		{Group: "b", User: "jack"},
	}
	
	ii := UniqueMap[I, string](items, func(item I) (key string, v I) {
		return item.Group, item
	})
	
	t.Log(ii)
}
