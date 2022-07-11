package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	. "github.com/lucas-s-work/funcy-go/iterator"
	"github.com/lucas-s-work/funcy-go/slice"
)

func main() {
	// Perform some maths and then string mapping on the below array
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	i1 := NewSliceIterator(a)
	i2 := Filter(i1, func(v int) (bool, error) {
		return v%2 == 0, nil
	})
	i2 = Map(i2, func(v int) (int, error) {
		return v * v, nil
	})
	i3 := Map(i2, func(v int) (int, error) {
		return v * 2, nil
	})
	i4 := Map(i3, func(v int) (string, error) { return fmt.Sprintf("value is: %v ", v), nil })
	fmt.Println(Sum(i4))

	// Do the same but using a natural number generator
	// Note the "WithLimit" to make the infinite generator finite
	n := WithLimit(NewNaturalGenerator(), 100)
	n1 := Filter(n, func(v int) (bool, error) {
		return v%2 == 0 && v != 0, nil
	})
	n2 := Map(n1, func(v int) (int, error) {
		return v * v, nil
	})
	n3 := Map(n2, func(v int) (int, error) {
		return v * 2, nil
	})
	n4 := Map(n3, func(v int) (string, error) { return fmt.Sprintf("value is: %v ", v), nil })
	fmt.Println(Sum(n4))

	// Multiply each consecutive pair of integers produced by i2
	err := i3.Reset()
	if err != nil {
		panic(err)
	}
	i5 := Cons(i3, 2)
	i6 := Map(i5, func(vs Iterator[int]) (int, error) {
		return Mult(vs)
	})

	c, err := Collect(i6)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)

	// Print every third fibonacci number
	f := NewFibonnacciGenerator()
	m := NewMaskGenerator(2)
	f2 := Filter(f, func(v int) (bool, error) {
		ok, err, _ := m.Next()

		return ok, err
	})
	fibs, _ := CollectWithLimit(f2, 10)
	fmt.Println(fibs)

	// Validation regex using an error
	amountRegex := regexp.MustCompile(`^[$£]?\d{1,13}\.?\d{0,13}$`)
	amounts := []string{"123.123", "£1.198", "$1", "not an amount", "99.9", "64"}
	ai1 := NewSliceIterator(amounts)
	if err := Each(ai1, func(am string) error {
		if !amountRegex.MatchString(am) {
			return fmt.Errorf("invalid amount: \"%s\"", am)
		}

		return nil
	}); err != nil {
		fmt.Println("Invalid amount found:", err)
	}

	// Exclusion and formatting to float of amounts
	ai1.Reset()
	ai2 := Filter(ai1, func(s string) (bool, error) {
		return amountRegex.MatchString(s), nil
	})
	ai3 := Map(ai2, func(s string) (string, error) {
		return strings.TrimPrefix(s, "£"), nil
	})
	ai4 := Map(ai3, func(s string) (string, error) {
		return strings.TrimPrefix(s, "$"), nil
	})
	ai5 := Map(ai4, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
	amountFloats, err := Collect(ai5)
	fmt.Println(amountFloats)

	// Use the slice package to perform useful operations on slices
	b := []float32{1, 5, 9, 2, 3, 1, 4}
	fmt.Println(b)
	b = slice.Sort(b)
	fmt.Println(b)
	fmt.Println(slice.Pick(b))
	fmt.Println(slice.Shuffle(b))
}
