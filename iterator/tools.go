package iterator

import (
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
