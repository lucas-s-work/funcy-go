package result

import "github.com/lucas-s-work/funcy-go/monad"

// An example of how a result monad might be implemented in golang
type Result[V any] struct {
	value V
	err   error
}

func Wrap[I any](i I) Result[I] {
	return Result[I]{
		value: i,
	}
}

func (r Result[V]) Unwrap() (V, error) {
	return r.value, r.err
}

func Bind[I, O any](r Result[I], f func(I) Result[O]) Result[O] {
	if err := r.err; err != nil {
		return Result[O]{err: err}
	}

	return f(r.value)
}

// Prove that result.Bind is a valid Bind function
var _ monad.MonadOb[int, error] = Result[int]{}
var _ monad.Binder[Result[int], Result[int], int, int, error] = Bind[int, int]
