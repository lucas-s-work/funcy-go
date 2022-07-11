package main

import (
	"fmt"
	"math"
	"testing"

	. "github.com/lucas-s-work/funcy-go/iterator"
)

func Expensive(f float64) (int, error) {
	i := int(f)
	i %= 10
	for j := 0; j < 50; j++ {
		i = int(math.Exp(float64(i))) % 10
	}

	return i, nil
}

// It's as fast as imperitive, and as memory efficient as imperitive most of the time too!
// go test -bench=. -benchmem ./...
// BenchmarkFunctional-8                 58          20375341 ns/op        39366147 B/op       8613 allocs/op
// BenchmarkImperitive-8                 54          20476251 ns/op        39366005 B/op       8606 allocs/op

func BenchmarkFunctional(b *testing.B) {
	// Take the first 10000 numbers, multiply them by 3, disgard multiples of 7 and convert to a string then concat
	for i := 0; i < b.N; i++ {
		ns := WithLimit(NewNaturalGenerator(), 10000)
		ns = Map(ns, func(i int) (int, error) {
			return i * 3, nil
		})
		ns = Filter(ns, func(i int) (bool, error) {
			return !(i%7 == 0), nil
		})
		ns2 := Map(ns, func(i int) (float64, error) {
			return math.Sqrt(float64(i)), nil
		})
		ns3 := Map(ns2, func(f float64) (float64, error) {
			return math.Sin(math.Cos(math.Sin(math.Tan(f)))), nil
		})
		ns4 := Map(ns3, Expensive)
		strs := Map(ns4, func(i int) (string, error) {
			return fmt.Sprintf("%d", i), nil
		})
		Sum(strs)
	}
}

func BenchmarkImperitive(b *testing.B) {
	// Perform the same but in normal golang
	for i := 0; i < b.N; i++ {
		str := ""
		for i := 0; i < 10000; i++ {
			v := i * 3
			if v%7 == 0 {
				continue
			}
			v2 := math.Sqrt(float64(v))
			v3 := math.Sin(math.Cos(math.Sin(math.Tan(v2))))
			v4, _ := Expensive(v3)
			str += fmt.Sprintf("%v", v4)
		}
	}
}
