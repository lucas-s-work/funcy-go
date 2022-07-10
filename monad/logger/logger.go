package logger

import (
	"strings"
)

// An example of a logger monad
type Logger[V any] struct {
	value V
	log   []string
}

func Wrap[V any](v V) Logger[V] {
	return Logger[V]{
		value: v,
		log:   make([]string, 0),
	}
}

func (l Logger[V]) Unwrap() (V, string) {
	return l.value, strings.Join(l.log, "\n")
}

func Bind[I, O any](l Logger[I], f func(I) Logger[O]) Logger[O] {
	nl := f(l.value)
	return Logger[O]{
		value: nl.value,
		log:   append(l.log, nl.log...),
	}
}
