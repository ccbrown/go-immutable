package immutable

// Stack implements a last in, first out container.
//
// Nil and the zero value for Stack are both empty stacks.
type Stack[T any] struct {
	top    T
	bottom *Stack[T]
}

// Empty returns true if the stack is empty.
//
// Complexity: O(1) worst-case
func (s *Stack[T]) Empty() bool {
	return s == nil || s.bottom == nil
}

// Peek returns the top item on the stack.
//
// Complexity: O(1) worst-case
func (s *Stack[T]) Peek() T {
	return s.top
}

// Pop removes the top item from the stack.
//
// Complexity: O(1) worst-case
func (s *Stack[T]) Pop() *Stack[T] {
	return s.bottom
}

// Push places an item onto the top of the stack.
//
// Complexity: O(1) worst-case
func (s *Stack[T]) Push(value T) *Stack[T] {
	if s == nil {
		s = &Stack[T]{}
	}
	return &Stack[T]{
		top:    value,
		bottom: s,
	}
}
