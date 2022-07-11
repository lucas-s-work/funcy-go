package iterator

import (
	"github.com/lucas-s-work/funcy-go/queue"
	"golang.org/x/exp/constraints"
)

func Each[V any](i Iterator[V], f func(V) error) error {
	for {
		v, err, ok := i.Next()
		if !ok {
			return nil
		}
		if err != nil {
			return err
		}

		if err := f(v); err != nil {
			return err
		}
	}
}

func EachWithLimit[V any](i Iterator[V], f func(V) error, limit int) error {
	for index := 0; index < limit; index++ {
		v, err, ok := i.Next()
		if !ok {
			return nil
		}
		if err != nil {
			return err
		}

		if err := f(v); err != nil {
			return err
		}
	}

	return nil
}

func Collect[V any](i Iterator[V]) ([]V, error) {
	var out []V
	if err := Each(i, func(v V) error {
		out = append(out, v)

		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func CollectWithLimit[V any](i Iterator[V], limit int) ([]V, error) {
	var out []V
	if err := EachWithLimit(i, func(v V) error {
		out = append(out, v)

		return nil
	}, limit); err != nil {
		return nil, err
	}

	return out, nil
}

type mappedIterator[I, O any] struct {
	Iterator[I]
	action func(I) (O, error)
}

func (m *mappedIterator[I, O]) Next() (O, error, bool) {
	v, err, ok := m.Iterator.Next()
	var o O
	if !ok {
		return o, nil, ok
	}
	if err != nil {
		return o, err, true
	}

	o, err = m.action(v)
	return o, err, true
}

func Map[I, O any](i Iterator[I], f func(I) (O, error)) Iterator[O] {
	return &mappedIterator[I, O]{
		Iterator: i,
		action:   f,
	}
}

type filterIterator[I any] struct {
	Iterator[I]
	check func(I) (bool, error)
}

func (f *filterIterator[I]) Next() (I, error, bool) {
	for {
		v, err, ok := f.Iterator.Next()
		if !ok {
			return v, nil, false
		}
		if err != nil {
			return v, err, true
		}

		ok, err = f.check(v)
		if ok {
			return v, nil, true
		}
		if err != nil {
			var o I
			return o, err, true
		}
	}
}

func Filter[I any](i Iterator[I], c func(I) (bool, error)) Iterator[I] {
	return &filterIterator[I]{
		Iterator: i,
		check:    c,
	}
}

type ConsIterator[V any, O Iterator[V]] struct {
	Iterator[V]
	stride int
	empty  bool
}

func Cons[I any](i Iterator[I], stride int) Iterator[Iterator[I]] {
	return &ConsIterator[I, Iterator[I]]{
		Iterator: i,
		stride:   stride,
		empty:    false,
	}
}

func (c *ConsIterator[I, O]) Next() (Iterator[I], error, bool) {
	if c.empty {
		return nil, nil, false
	}
	var out []I
	for i := 0; i < c.stride; i++ {
		v, err, ok := c.Iterator.Next()
		if !ok {
			c.empty = true
			if len(out) == 0 {
				return nil, nil, false
			}

			return NewSliceIterator(out), nil, true
		}
		if err != nil {
			return nil, err, true
		}

		out = append(out, v)
	}

	return NewSliceIterator(out), nil, true
}

func (c *ConsIterator[I, O]) Reset() error {
	c.empty = false

	return c.Iterator.Reset()
}

func Fold[I, O any](i Iterator[I], acc O, f func(I, O) (O, error)) (O, error) {
	if err := Each(i, func(v I) error {
		var err error
		acc, err = f(v, acc)

		return err
	}); err != nil {
		var o O
		return o, err
	}

	return acc, nil
}

type scanIterator[I, O any] struct {
	in  Iterator[I]
	f   func(I, O) (O, error)
	acc O
}

func Scan[I, O any](i Iterator[I], acc O, f func(I, O) (O, error)) Iterator[O] {
	return &scanIterator[I, O]{
		in:  i,
		f:   f,
		acc: acc,
	}
}

func (s *scanIterator[I, O]) Next() (O, error, bool) {
	v, err, ok := s.in.Next()
	var o O
	if !ok {
		return o, nil, false
	}
	if err != nil {
		return o, err, true
	}

	s.acc, err = s.f(v, s.acc)
	if err != nil {
		return o, err, true
	}
	return s.acc, nil, true
}

func (s *scanIterator[I, O]) Reset() error {
	if err := s.in.Reset(); err != nil {
		return err
	}

	var zero O
	s.acc = zero
	return nil
}

func Sum[V constraints.Ordered](i Iterator[V]) (V, error) {
	var acc V
	return Fold(i, acc, func(v, acc V) (V, error) { return acc + v, nil })
}

type Multable interface {
	constraints.Complex | constraints.Float | constraints.Integer
}

func Mult[V Multable](i Iterator[V]) (V, error) {
	return Fold(i, 1, func(v, acc V) (V, error) { return acc * v, nil })
}

type distinctIterator[V any, C constraints.Ordered] struct {
	in    Iterator[V]
	seen  map[C]struct{}
	check func(V) C
}

func Distinct[V constraints.Ordered](i Iterator[V]) Iterator[V] {
	return &distinctIterator[V, V]{
		in:   i,
		seen: make(map[V]struct{}),
		check: func(v V) V {
			return v
		},
	}
}

func DistinctMap[V any, C constraints.Ordered](i Iterator[V], c func(V) C) Iterator[V] {
	return &distinctIterator[V, C]{
		in:    i,
		seen:  make(map[C]struct{}),
		check: c,
	}
}

func (d *distinctIterator[V, C]) Next() (V, error, bool) {
	for {
		v, err, ok := d.in.Next()
		if !ok {
			return v, nil, false
		}
		if err != nil {
			return v, err, true
		}

		c := d.check(v)
		if _, seen := d.seen[c]; !seen {
			d.seen[c] = struct{}{}
			return v, nil, true
		}
	}
}

func (d *distinctIterator[V, C]) Reset() error {
	err := d.in.Reset()
	if err != nil {
		return err
	}

	d.seen = make(map[C]struct{})
	return nil
}

type mergeIterator[V any] struct {
	in1, in2 Iterator[V]
	swap     bool
}

func Merge[I any](i1, i2 Iterator[I]) Iterator[I] {
	return &mergeIterator[I]{
		in1:  i1,
		in2:  i2,
		swap: false,
	}
}

func (m *mergeIterator[V]) Next() (V, error, bool) {
	if m.swap {
		v, err, ok := m.in2.Next()
		if !ok {
			return m.in1.Next()
		}
		if err != nil {
			return v, err, true
		}

		m.swap = !m.swap
		return v, err, ok
	}

	v, err, ok := m.in1.Next()
	if !ok {
		return m.in2.Next()
	}
	if err != nil {
		return v, err, true
	}

	m.swap = !m.swap
	return v, err, ok
}

func (m *mergeIterator[V]) Reset() error {
	if err := m.in1.Reset(); err != nil {
		return err
	}
	if err := m.in2.Reset(); err != nil {
		return err
	}

	m.swap = false
	return nil
}

// Alternatively use fold
func Max[V constraints.Ordered](i Iterator[V]) (V, error) {
	var m V
	set := false

	if err := Each(i, func(v V) error {
		if !set || v > m {
			set = true
			m = v
		}

		return nil
	}); err != nil {
		return m, err
	}

	return m, nil
}

func Min[V constraints.Ordered](i Iterator[V]) (V, error) {
	var m V
	set := false

	if err := Each(i, func(v V) error {
		if !set || v < m {
			set = true
			m = v
		}

		return nil
	}); err != nil {
		return m, err
	}

	return m, nil
}

type split[V any] struct {
	i       Iterator[V]
	check   func(V) bool
	partner *split[V]
	cache   queue.Queue[V]
}

func (s *split[V]) push(v V) {
	s.cache.Push(v)
}

func Partition[V any](i Iterator[V], check func(V) bool) (Iterator[V], Iterator[V]) {
	t := &split[V]{
		i:     i,
		check: check,
		cache: queue.Queue[V]{},
	}
	f := &split[V]{
		i:     i,
		check: func(v V) bool { return !check(v) },
		cache: queue.Queue[V]{},
	}
	t.partner = f
	f.partner = t

	return t, f
}

func (s *split[V]) Next() (V, error, bool) {
	// Pull off the cache first
	v, ok := s.cache.Pop()
	if ok {
		return v, nil, true
	}

	// If the cache is empty iterate until we get a value or none is found
	for {
		v, err, ok := s.i.Next()
		if !ok {
			return v, nil, false
		}
		if err != nil {
			return v, err, true
		}

		// Check if v belongs to this iterator, if not place it on our partners cache and keep going
		if s.check(v) {
			return v, nil, true
		}
		s.partner.push(v)
	}
}

func (s *split[V]) Reset() error {
	if err := s.i.Reset(); err != nil {
		return err
	}

	s.cache = queue.Queue[V]{}
	s.partner.cache = queue.Queue[V]{}

	return nil
}

func GroupBy[K constraints.Ordered, V any](i Iterator[V], f func(V) K) (map[K][]V, error) {
	return Fold(i, make(map[K][]V), func(v V, m map[K][]V) (map[K][]V, error) {
		k := f(v)
		arr, ok := m[k]
		if !ok {
			m[k] = []V{v}

			return m, nil
		}

		m[k] = append(arr, v)
		return m, nil
	})
}
