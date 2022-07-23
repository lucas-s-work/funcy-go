package queue

import "fmt"

const indexLim = 100

type Queue[V any] struct {
	start, end *Node[V]
}

type Node[V any] struct {
	el   V
	prev *Node[V]
}

func (q *Queue[V]) Push(v V) {
	n := &Node[V]{
		el: v,
	}
	if q.start == nil {
		q.start = n
		q.end = q.start

		return
	}

	q.start.prev = n
	q.start = n
}

func (q *Queue[V]) Pop() (V, bool) {
	if q.end == nil {
		var o V
		return o, false
	}

	end := q.end
	if q.end == q.start {
		q.end = nil
		q.start = nil
	} else {
		q.end = end.prev
	}

	return end.el, true
}

// Implement the iterator interface
func (q *Queue[V]) Next() (V, error, bool) {
	v, ok := q.Pop()
	return v, nil, ok
}

func (q *Queue[V]) Reset() error {
	return fmt.Errorf("cannot reset queue type")
}
