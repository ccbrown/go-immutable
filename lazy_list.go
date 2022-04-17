package immutable

import (
	"sync"
)

type lazyList[T any] struct {
	value      T
	lazyNext   func() *lazyList[T]
	next       *lazyList[T]
	evaluation sync.Once
}

func newLazyList[T any](front T, next func() *lazyList[T]) *lazyList[T] {
	return &lazyList[T]{
		value:    front,
		lazyNext: next,
	}
}

func (l *lazyList[T]) Front() T {
	return l.value
}

func (l *lazyList[T]) PopFront() *lazyList[T] {
	l.evaluation.Do(func() {
		if l.lazyNext != nil {
			l.next = l.lazyNext()
			l.lazyNext = nil
		}
	})
	return l.next
}

func (l *lazyList[T]) PushFront(value T) *lazyList[T] {
	return &lazyList[T]{
		value: value,
		next:  l,
	}
}
