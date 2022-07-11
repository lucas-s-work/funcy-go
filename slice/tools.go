package slice

import (
	"math/rand"
	"sort"
	"time"

	"golang.org/x/exp/constraints"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Pick[V any](s []V) V {
	var o V
	if len(s) == 0 {
		return o
	}
	return s[rand.Int()%len(s)]
}

func Shuffle[V any](s []V) []V {
	if len(s) == 0 {
		return s
	}

	m := make(map[int]struct{})
	for i := range s {
		m[i] = struct{}{}
	}

	j := 0
	o := make([]V, len(s))
	for i := range m {
		o[j] = s[i]
		j++
	}

	return o
}

func Reverse[V constraints.Ordered](arr []V) []V {
	out := make([]V, len(arr))
	l := len(out)
	for i := range arr {
		out[l-i-1] = arr[i]
	}

	return out
}

type sortable[V constraints.Ordered] []V

func (s sortable[V]) Len() int {
	return len(s)
}

func (s sortable[V]) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortable[V]) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func Sort[V constraints.Ordered](s []V) []V {
	var ss sortable[V] = s
	sort.Sort(ss)
	var o []V = ss
	return o
}
