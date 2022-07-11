package stack

type Stack[V any] struct {
	top *Node[V]
}

type Node[V any] struct {
	el   V
	prev *Node[V]
}

func (s *Stack[V]) Pop() (V, bool) {
	top := s.top
	if top == nil {
		var o V
		return o, false
	}

	s.top = top.prev
	return top.el, true
}

func (s *Stack[V]) Push(v V) {
	if s.top == nil {
		s.top = &Node[V]{
			el: v,
		}
		return
	}

	top := s.top
	s.top = &Node[V]{
		el:   v,
		prev: top,
	}
}

// Implement the iterator interface
func (s *Stack[V]) Next() (V, error, bool) {
	v, ok := s.Pop()
	return v, nil, ok
}
