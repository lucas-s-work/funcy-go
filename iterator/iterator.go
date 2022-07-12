package iterator

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// A set of Functors from Type -> Iterator[Type]
type Iterator[V any] interface {
	Next() (V, error, bool)
	Reset() error
}

type SliceIterator[V any] struct {
	s     []V
	index int
}

func NewSliceIterator[V any](s []V) Iterator[V] {
	return &SliceIterator[V]{
		s:     s,
		index: 0,
	}
}

func (s *SliceIterator[V]) Next() (V, error, bool) {
	if s.index >= len(s.s) {
		var v V
		return v, nil, false
	}

	s.index++
	return s.s[s.index-1], nil, true
}

func (s *SliceIterator[V]) Reset() error {
	s.index = 0
	return nil
}

var _ Iterator[int] = &SliceIterator[int]{}

type MapIterator[K constraints.Ordered, V any, O KeyValue[K, V]] struct {
	m        map[K]V
	keys     []K
	keyIndex int
}

type KeyValue[K constraints.Ordered, V any] struct {
	Key   K
	Value V
}

func NewMapIterator[K constraints.Ordered, V any](m map[K]V) Iterator[KeyValue[K, V]] {
	keys := make([]K, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
	}

	return &MapIterator[K, V, KeyValue[K, V]]{
		m:        m,
		keys:     keys,
		keyIndex: 0,
	}
}

func (m *MapIterator[K, V, O]) Next() (O, error, bool) {
	if m.keyIndex > len(m.keys) {
		var o O
		return o, nil, false
	}

	k := m.keys[m.keyIndex]
	o := O{
		Key:   k,
		Value: m.m[k],
	}

	m.keyIndex++
	return o, nil, true
}

func (s *MapIterator[K, V, O]) Reset() error {
	s.keyIndex = 0
	return nil
}

var _ Iterator[KeyValue[int, int]] = &MapIterator[int, int]{}

type ChanIterator[V any] struct {
	c    chan V
	outs []V
}

func NewChanIterator[V any](c chan V) Iterator[V] {
	return &ChanIterator[V]{
		c:    c,
		outs: []V{},
	}
}

func (c *ChanIterator[V]) Next() (V, error, bool) {
	v, ok := <-c.c
	return v, nil, ok
}

func (s *ChanIterator[V]) Reset() error {
	var v V
	return fmt.Errorf("cannot reset channel iterator for type: %T", v)
}

var _ Iterator[int] = &ChanIterator[int]{}

type MatrixIterator[V any] struct {
	m         [][]V
	dim1index int
	dim2index int
}

func NewMatrixIterator[V any](m [][]V) Iterator[V] {
	return &MatrixIterator[V]{
		m:         m,
		dim1index: 0,
		dim2index: 0,
	}
}

func (m *MatrixIterator[V]) Next() (V, error, bool) {
	if m.dim2index >= len(m.m[m.dim1index]) {
		m.dim2index = 0
		m.dim1index++
	}

	if m.dim1index >= len(m.m) {
		var v V
		return v, nil, false
	}

	m.dim2index++
	return m.m[m.dim1index][m.dim2index], nil, true
}

func (s *MatrixIterator[V]) Reset() error {
	s.dim1index = 0
	s.dim2index = 0
	return nil
}

var _ Iterator[int] = &MatrixIterator[int]{}

type StringIterator struct {
	s     string
	index int
}

func NewStringIterator(s string) Iterator[string] {
	return &StringIterator{
		s:     s,
		index: 0,
	}
}

func (s *StringIterator) Next() (string, error, bool) {
	if s.index >= len(s.s) {
		return "", nil, false
	}

	s.index++
	return string(s.s[s.index-1]), nil, true
}

func (s *StringIterator) Reset() error {
	s.index = 0
	return nil
}

var _ Iterator[int] = &SliceIterator[int]{}

type LimitedIterator[V any] struct {
	Iterator[V]
	limit int
	index int
}

func WithLimit[V any](i Iterator[V], limit int) Iterator[V] {
	return &LimitedIterator[V]{
		Iterator: i,
		limit:    limit,
		index:    0,
	}
}

func (l *LimitedIterator[V]) Next() (V, error, bool) {
	v, err, ok := l.Iterator.Next()
	if !ok {
		return v, err, false
	}
	if err != nil {
		return v, err, true
	}

	if l.index == l.limit {
		var o V
		return o, nil, false
	}
	l.index++

	return v, nil, true
}

func (l *LimitedIterator[V]) Reset() error {
	if err := l.Iterator.Reset(); err != nil {
		return err
	}
	l.index = 0

	return nil
}
