package unique

import "testing"

func TestSlices(t *testing.T) {
	a := []string{"a", "b", "c", "d", "d", "c", "a"}
	b := SliceKeys[string, string](a, func(item string) (key string) {
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
	
	ii := SliceKeys[I, string](items, func(item I) (key string) {
		return item.Group
	})
	
	t.Log(ii)
}
