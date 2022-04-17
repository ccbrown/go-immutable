package immutable

func queueRotate[T any](f *lazyList[T], r *Stack[T], s *lazyList[T]) *lazyList[T] {
	if f == nil {
		return s.PushFront(r.Peek())
	}
	return newLazyList(f.Front(), func() *lazyList[T] {
		return queueRotate(f.PopFront(), r.Pop(), s.PushFront(r.Peek()))
	})
}

func queueExec[T any](f *lazyList[T], r *Stack[T], s *lazyList[T]) *Queue[T] {
	if s == nil {
		f2 := queueRotate(f, r, nil)
		return &Queue[T]{f2, nil, f2}
	}
	return &Queue[T]{f, r, s.PopFront()}
}

// Queue implements a first in, first out container.
//
// Nil and the zero value for Queue are both empty queues.
type Queue[T any] struct {
	f *lazyList[T]
	r *Stack[T]
	s *lazyList[T]
}

// Empty returns true if the queue is empty.
//
// Complexity: O(1) worst-case
func (q *Queue[T]) Empty() bool {
	return q == nil || q.f == nil
}

// Front returns the item at the front of the queue.
//
// Complexity: O(1) worst-case
func (q *Queue[T]) Front() T {
	return q.f.Front()
}

// PopFront removes the item at the front of the queue.
//
// Complexity: O(1) worst-case
func (q *Queue[T]) PopFront() *Queue[T] {
	return queueExec(q.f.PopFront(), q.r, q.s)
}

// PushBack pushes an item onto the back of the queue.
//
// Complexity: O(1) worst-case
func (q *Queue[T]) PushBack(value T) *Queue[T] {
	return queueExec(q.f, q.r.Push(value), q.s)
}
