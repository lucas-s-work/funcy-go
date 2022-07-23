package iterator

import (
	"fmt"
	"math/rand"
)

type NaturalGenerator struct {
	i int
}

func NewNaturalGenerator() Iterator[int] {
	return &NaturalGenerator{}
}

func (n *NaturalGenerator) Next() (int, error, bool) {
	n.i++
	return n.i - 1, nil, true
}

func (n *NaturalGenerator) Reset() error {
	n.i = 0
	return nil
}

type FibonnacciGenerator struct {
	a, b int
}

func NewFibonnacciGenerator() Iterator[int] {
	return &FibonnacciGenerator{
		a: 1,
		b: 0,
	}
}

func (f *FibonnacciGenerator) Next() (int, error, bool) {
	c := f.a
	f.a = f.a + f.b
	f.b = c

	return c, nil, true
}

func (f *FibonnacciGenerator) Reset() error {
	f.a = 1
	f.b = 0

	return nil
}

type MaskGenerator struct {
	stride int
	index  int
}

func NewMaskGenerator(stride int) Iterator[bool] {
	return &MaskGenerator{
		stride: stride,
		index:  stride,
	}
}

func (m *MaskGenerator) Next() (bool, error, bool) {
	pass := m.stride == m.index

	m.index++
	if pass {
		m.index = 0
	}

	return pass, nil, true
}

func (m *MaskGenerator) Reset() error {
	m.index = m.stride
	return nil
}

// computed primes are globally known
var primes = []int{2, 3, 5}

type PrimeGenerator struct {
	index int
}

// Currently awfully slow, a better method may be to have an expanding memoized hashmap
// using a sieve
func NewPrimeGenerator() Iterator[int] {
	return &PrimeGenerator{
		index: 0,
	}
}

func (p *PrimeGenerator) Next() (int, error, bool) {
	// Don't re-compute the primes even if this is reset to save speed :)
	if p.index < len(primes) {
		p.index++
		return primes[p.index-1], nil, true
	}

	next := primes[p.index-1] + 2 // primes are spaced by atleast 2
	n2 := next / 2

	for {
		// no prime > 2 divisible by 2
		for _, prime := range primes {
			if prime > n2+1 {
				// We've found a new prime
				primes = append(primes, next)
				p.index++
				return next, nil, true
			}

			if next%prime == 0 {
				break
			}
		}
		// No prime found, try next number
		next += 2
		n2 = next / 2
	}
}

func (p *PrimeGenerator) Reset() error {
	p.index = 0

	return nil
}

type GeneratorFromFunction[V any] func() V

func FromFunction[V any](f func() V) Iterator[V] {
	return GeneratorFromFunction[V](f)
}

func (g GeneratorFromFunction[V]) Next() (V, error, bool) {
	return g(), nil, true
}

func (g GeneratorFromFunction[V]) Reset() error {
	return fmt.Errorf("cannot reset generating function")
}

func NewRandomIntGenerator() Iterator[int] {
	return GeneratorFromFunction[int](rand.Int)
}

func NewRandomFloat32Generator() Iterator[float32] {
	return GeneratorFromFunction[float32](rand.Float32)
}

func NewRandomFloat64Generator() Iterator[float64] {
	return GeneratorFromFunction[float64](rand.Float64)
}
