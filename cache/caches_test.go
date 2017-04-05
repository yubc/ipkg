package cache

import (
	"fmt"
	"sort"
	"strconv"
	"testing"
)

func tCache(key string) (interface{}, error) {
	var a = []int{1, 2, 5, 3, 4, 6, 8, 7, 9, 10}
	sort.Sort(sort.IntSlice(a))
	return a, nil
}

func TestCache(t *testing.T) {
	var m *Memo
	m = New(tCache)
	fmt.Println(m.Get("100"))
}

func BenchmarkCache(b *testing.B) {
	var m *Memo

	m = New(tCache)
	for i := 0; i < b.N; i++ {
		fmt.Println(m.Get(strconv.Itoa(i)))
	}
}
