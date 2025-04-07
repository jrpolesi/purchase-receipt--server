package testutils

import (
	"go.uber.org/mock/gomock"
)

func MatchedByFn[T any](fn func(T) bool) gomock.Matcher {
	return &matcher[T]{fn: fn}
}

type matcher[T any] struct {
	fn func(T) bool
}

func (m *matcher[T]) Matches(x any) bool {
	val, ok := x.(T)
	return ok && m.fn(val)
}

func (m *matcher[T]) String() string {
	return "matcher for custom condition"
}
