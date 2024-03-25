package sorts

import "testing"

type User struct {
	Id   int
	Name string
}

func TestSort(t *testing.T) {
	list := []int{1, 3, 4, 2, 5}
	strs := []string{"a", "c", "d", "e", "b"}

	Sort(list, func(a, b int) bool {
		return a > b
	})
	t.Log(list)

	Sort(strs, func(a, b string) bool {
		return a < b
	})
	t.Log(strs)
}

func TestSortAsc(t *testing.T) {
	list := []int{1, 3, 4, 2, 5}
	SortAsc(list, func(item int) int {
		return item
	})
	t.Log(list)

	users := []User{
		{Id: 1, Name: "tom"},
		{Id: 2, Name: "jack"},
		{Id: 3, Name: "jack"},
	}
	SortAsc(users, func(item User) int {
		return item.Id
	})
	t.Log(users)
	SortDesc(users, func(item User) int {
		return item.Id
	})
	t.Log(users)
}
