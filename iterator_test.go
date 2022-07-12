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
	for j := 0; j < 10; j++ {
		i = int(math.Exp(float64(i))) % 10
	}

	return i, nil
}

// For significant workloads its difference to imperitive is negligible, and its allocations are effectively identical
// go test -bench=. -benchmem ./...
// BenchmarkFunctional-8                 28          45235673 ns/op        393242076 B/op     17472 allocs/op
// BenchmarkImperitive-8                 25          45932572 ns/op        393241758 B/op     17463 allocs/op

type output struct {
	message string
	values  []int
}

func BenchmarkFunctional(b *testing.B) {
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
		out := output{
			message: "",
			values:  []int{},
		}
		out, _ = Fold(ns4, out, func(i int, out output) (output, error) {
			out.message += fmt.Sprintf("value is %v", i)
			out.values = append(out.values, i)
			return out, nil
		})
	}
}

func BenchmarkImperitive(b *testing.B) {
	// Perform the same but in normal golang
	for i := 0; i < b.N; i++ {
		out := output{
			message: "",
			values:  []int{},
		}
		for j := 0; j < 10000; j++ {
			v := j * 3
			if v%7 == 0 {
				continue
			}
			v2 := math.Sqrt(float64(v))
			v3 := math.Sin(math.Cos(math.Sin(math.Tan(v2))))
			v4, _ := Expensive(v3)
			str := fmt.Sprintf("value is %v", v4)
			out.message += str
			out.values = append(out.values, v4)
		}
	}
}

const iterations = 100

// For very simple workloads it's slower by about 3x but the memory allocation is the same
// For comparison with another golang library it is still 2x as fast and uses half the memory
// Mine:
// BenchmarkImperative2-8           4291770               276.6 ns/op          1016 B/op          7 allocs/op
// BenchmarkFunctional2-8           1412942               846.2 ns/op          1104 B/op         11 allocs/op
// Theirs:
// BenchmarkImperative-4            2098518               550.6 ns/op          1016 B/op          7 allocs/op
// BenchmarkFunctional-4             293095              3653 ns/op            2440 B/op         23 allocs/op

func BenchmarkImperative2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		var result []int
		for i := 0; i < iterations; i++ {
			if count%3 == 0 {
				result = append(result, count*count)
			}
			count++
		}
		_ = result
	}
}

func BenchmarkFunctional2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		n := WithLimit(NewNaturalGenerator(), iterations)
		n = Filter(n, func(i int) (bool, error) {
			return i%3 == 0, nil
		})
		n = Map(n, func(i int) (int, error) {
			return i * i, nil
		})
		Collect(n)
	}
}

func BenchmarkComplexImperitive(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		for i := 0; i < 5000; i++ {
			if (i*i*7)%17 != 0 {
				continue
			}

			v := math.Cos(float64(i)) * float64(i)
			s := fmt.Sprintf("value is: %v", v)
			if s[0:2] == "val" {
				count++
			}
		}
	}
}

func BenchmarkComplexFunctional(b *testing.B) {
	for n := 0; n < b.N; n++ {
		n := WithLimit(NewNaturalGenerator(), 5000)
		n = Filter(n, func(i int) (bool, error) {
			return int(i*i*7)%17 == 0, nil
		})
		n2 := Map(n, func(i int) (float64, error) {
			return math.Cos(float64(i)) * float64(i), nil
		})
		n3 := Map(n2, func(i float64) (string, error) {
			return fmt.Sprintf("value is: %v", i), nil
		})
		CountFilter(n3, func(i string) bool {
			return i[0:2] == "val"
		})
	}
}
