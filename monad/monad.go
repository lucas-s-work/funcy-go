package monad

// What a monad signature might look like in golang

type MonadOb[V, S any] interface {
	Unwrap() (V, S)
}

type Binder[IM MonadOb[I, S], OM MonadOb[O, S], I, O, S any] func(IM, func(I) OM) OM
