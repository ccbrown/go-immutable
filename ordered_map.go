package immutable

import "golang.org/x/exp/constraints"

const (
	orderedMapNegativeBlack = -1
	orderedMapRed           = 0
	orderedMapBlack         = 1
	orderedMapDoubleBlack   = 2
)

// OrderedMap implements an ordered map.
//
// Nil and the zero value for OrderedMap are both empty maps.
type OrderedMap[K constraints.Ordered, V any] struct {
	len   int
	color int
	left  *OrderedMap[K, V]
	right *OrderedMap[K, V]
	key   K
	value V
}

// Empty returns true if the map is empty.
//
// Complexity: O(1) worst-case
func (m *OrderedMap[K, V]) Empty() bool {
	return m == nil || m.len == 0
}

// Len returns the number of elements in the map.
//
// Complexity: O(1) worst-case
func (m *OrderedMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return m.len
}

// Get returns the value associated with the given key if set.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) Get(key K) (v V, exists bool) {
	l := m.findLessThanOrEqual(key, nil)
	if l == nil {
		return v, false
	}
	if l.key >= key {
		return l.value, true
	}
	return v, false
}

// Set associates a value with the given key.
//
// Only the built-in types may be used as keys. Once a value is set within a map, all subsequent
// operations must use the same key type.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) Set(key K, value V) *OrderedMap[K, V] {
	ret := m.insert(key, value)
	ret.color = orderedMapBlack
	return ret
}

// Delete removes a key from the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) Delete(key K) *OrderedMap[K, V] {
	if ret, _ := m.delete(key); !ret.Empty() {
		ret.color = orderedMapBlack
		return ret
	}
	return nil
}

// Min returns the minimum element in the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) Min() *OrderedMapElement[K, V] {
	return m.min(nil)
}

// Max returns the maximum element in the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) Max() *OrderedMapElement[K, V] {
	return m.max(nil)
}

// MinAfter returns the minimum element in the map that is greater than the given key.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) MinAfter(key K) *OrderedMapElement[K, V] {
	return m.minGreaterThan(key, nil)
}

// MaxBefore returns the maximum element in the map that is less than the given key.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap[K, V]) MaxBefore(key K) *OrderedMapElement[K, V] {
	return m.maxLessThan(key, nil)
}

func (m *OrderedMap[K, V]) min(lineage *Stack[*OrderedMap[K, V]]) *OrderedMapElement[K, V] {
	if m.Empty() {
		return nil
	} else if m.left != nil {
		return m.left.min(lineage.Push(m))
	}
	return &OrderedMapElement[K, V]{
		lineage: lineage,
		element: m,
	}
}

func (m *OrderedMap[K, V]) max(lineage *Stack[*OrderedMap[K, V]]) *OrderedMapElement[K, V] {
	if m.Empty() {
		return nil
	} else if m.right != nil {
		return m.right.max(lineage.Push(m))
	}
	return &OrderedMapElement[K, V]{
		lineage: lineage,
		element: m,
	}
}

func (m *OrderedMap[K, V]) minGreaterThan(key K, lineage *Stack[*OrderedMap[K, V]]) *OrderedMapElement[K, V] {
	if m.Empty() {
		return nil
	} else if key < m.key {
		if m.left != nil {
			if r := m.left.minGreaterThan(key, lineage.Push(m)); r != nil {
				return r
			}
		}
		return &OrderedMapElement[K, V]{
			lineage: lineage,
			element: m,
		}
	} else if m.key < key {
		return m.right.minGreaterThan(key, lineage.Push(m))
	}
	return m.right.min(lineage.Push(m))
}

func (m *OrderedMap[K, V]) maxLessThan(key K, lineage *Stack[*OrderedMap[K, V]]) *OrderedMapElement[K, V] {
	if m.Empty() {
		return nil
	} else if m.key < key {
		if m.right != nil {
			if r := m.right.maxLessThan(key, lineage.Push(m)); r != nil {
				return r
			}
		}
		return &OrderedMapElement[K, V]{
			lineage: lineage,
			element: m,
		}
	} else if key < m.key {
		return m.left.maxLessThan(key, lineage.Push(m))
	}
	return m.left.max(lineage.Push(m))
}

func (m *OrderedMap[K, V]) delete(key K) (*OrderedMap[K, V], bool) {
	if m.Empty() {
		return m, false
	} else if key < m.key {
		if left, didDelete := m.left.delete(key); didDelete {
			return m.adopt(left, m.right).bubble(), true
		}
		return m, false
	} else if m.key < key {
		if right, didDelete := m.right.delete(key); didDelete {
			return m.adopt(m.left, right).bubble(), true
		}
		return m, false
	}
	return m.remove(), true
}

func (m *OrderedMap[K, V]) adopt(left, right *OrderedMap[K, V]) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		len:   1 + left.Len() + right.Len(),
		color: m.color,
		left:  left,
		right: right,
		key:   m.key,
		value: m.value,
	}
}

func (m *OrderedMap[K, V]) findLessThanOrEqual(key K, candidate *OrderedMap[K, V]) *OrderedMap[K, V] {
	if m.Empty() {
		return candidate
	} else if key < m.key {
		return m.left.findLessThanOrEqual(key, candidate)
	}
	return m.right.findLessThanOrEqual(key, m)
}

func (m *OrderedMap[K, V]) insert(key K, value V) *OrderedMap[K, V] {
	if m.Empty() {
		return &OrderedMap[K, V]{
			len:   1,
			color: orderedMapRed,
			key:   key,
			value: value,
		}
	} else if key < m.key {
		return m.adopt(m.left.insert(key, value), m.right).balanceLeft()
	} else if m.key < key {
		return m.adopt(m.left, m.right.insert(key, value)).balanceRight()
	}
	return &OrderedMap[K, V]{
		len:   m.len,
		color: m.color,
		left:  m.left,
		right: m.right,
		key:   m.key,
		value: value,
	}
}

func (m *OrderedMap[K, V]) balanceLeft() *OrderedMap[K, V] {
	if m.color >= orderedMapBlack && m.left != nil {
		if m.left.color == orderedMapRed {
			if m.left.left != nil && m.left.left.color == orderedMapRed {
				return &OrderedMap[K, V]{
					len:   m.len,
					color: m.color - 1,
					left: &OrderedMap[K, V]{
						len:   m.left.left.len,
						color: orderedMapBlack,
						left:  m.left.left.left,
						right: m.left.left.right,
						key:   m.left.left.key,
						value: m.left.left.value,
					},
					right: &OrderedMap[K, V]{
						len:   1 + m.left.right.Len() + m.right.Len(),
						color: orderedMapBlack,
						left:  m.left.right,
						right: m.right,
						key:   m.key,
						value: m.value,
					},
					key:   m.left.key,
					value: m.left.value,
				}
			} else if m.left.right != nil && m.left.right.color == orderedMapRed {
				return &OrderedMap[K, V]{
					len:   m.len,
					color: m.color - 1,
					left: &OrderedMap[K, V]{
						len:   1 + m.left.left.Len() + m.left.right.left.Len(),
						color: orderedMapBlack,
						left:  m.left.left,
						right: m.left.right.left,
						key:   m.left.key,
						value: m.left.value,
					},
					right: &OrderedMap[K, V]{
						len:   1 + m.left.right.right.Len() + m.right.Len(),
						color: orderedMapBlack,
						left:  m.left.right.right,
						right: m.right,
						key:   m.key,
						value: m.value,
					},
					key:   m.left.right.key,
					value: m.left.right.value,
				}
			}
		} else if m.left.color == orderedMapNegativeBlack {
			left := &OrderedMap[K, V]{
				len:   1 + m.left.left.Len() + m.left.right.left.Len(),
				color: orderedMapBlack,
				left:  m.left.left.redden(),
				right: m.left.right.left,
				key:   m.left.key,
				value: m.left.value,
			}
			left = left.balanceLeft()
			right := &OrderedMap[K, V]{
				len:   1 + m.left.right.right.Len() + m.right.Len(),
				color: orderedMapBlack,
				left:  m.left.right.right,
				right: m.right,
				key:   m.key,
				value: m.value,
			}
			return &OrderedMap[K, V]{
				len:   1 + left.Len() + right.Len(),
				color: orderedMapBlack,
				left:  left,
				right: right,
				key:   m.left.right.key,
				value: m.left.right.value,
			}
		}
	}
	return m
}

func (m *OrderedMap[K, V]) balanceRight() *OrderedMap[K, V] {
	if m.color >= orderedMapBlack && m.right != nil {
		if m.right.color == orderedMapRed {
			if m.right.left != nil && m.right.left.color == orderedMapRed {
				return &OrderedMap[K, V]{
					len:   m.len,
					color: m.color - 1,
					left: &OrderedMap[K, V]{
						len:   1 + m.left.Len() + m.right.left.left.Len(),
						color: orderedMapBlack,
						left:  m.left,
						right: m.right.left.left,
						key:   m.key,
						value: m.value,
					},
					right: &OrderedMap[K, V]{
						len:   1 + m.right.left.right.Len() + m.right.right.Len(),
						color: orderedMapBlack,
						left:  m.right.left.right,
						right: m.right.right,
						key:   m.right.key,
						value: m.right.value,
					},
					key:   m.right.left.key,
					value: m.right.left.value,
				}
			} else if m.right.right != nil && m.right.right.color == orderedMapRed {
				return &OrderedMap[K, V]{
					len:   m.len,
					color: m.color - 1,
					left: &OrderedMap[K, V]{
						len:   1 + m.left.Len() + m.right.left.Len(),
						color: orderedMapBlack,
						left:  m.left,
						right: m.right.left,
						key:   m.key,
						value: m.value,
					},
					right: &OrderedMap[K, V]{
						len:   m.right.right.len,
						color: orderedMapBlack,
						left:  m.right.right.left,
						right: m.right.right.right,
						key:   m.right.right.key,
						value: m.right.right.value,
					},
					key:   m.right.key,
					value: m.right.value,
				}
			}
		} else if m.right.color == orderedMapNegativeBlack {
			left := &OrderedMap[K, V]{
				len:   1 + m.left.Len() + m.right.left.left.Len(),
				color: orderedMapBlack,
				left:  m.left,
				right: m.right.left.left,
				key:   m.key,
				value: m.value,
			}
			right := &OrderedMap[K, V]{
				len:   1 + m.right.left.right.Len() + m.right.right.Len(),
				color: orderedMapBlack,
				left:  m.right.left.right,
				right: m.right.right.redden(),
				key:   m.right.key,
				value: m.right.value,
			}
			right = right.balanceRight()
			return &OrderedMap[K, V]{
				len:   1 + left.Len() + right.Len(),
				color: orderedMapBlack,
				left:  left,
				right: right,
				key:   m.right.left.key,
				value: m.right.left.value,
			}
		}
	}
	return m
}

func (m *OrderedMap[K, V]) remove() *OrderedMap[K, V] {
	if !m.left.Empty() && !m.right.Empty() {
		left, removed := m.left.removeMax()
		reduced := &OrderedMap[K, V]{
			len:   m.len - 1,
			color: m.color,
			left:  left,
			right: m.right,
			key:   removed.key,
			value: removed.value,
		}
		return reduced.bubble()
	}
	var child *OrderedMap[K, V]
	if !m.left.Empty() {
		child = m.left
	} else if !m.right.Empty() {
		child = m.right
	} else {
		if m.color == orderedMapRed {
			return nil
		}
		return &OrderedMap[K, V]{color: orderedMapDoubleBlack}
	}
	ret := *child
	ret.color = orderedMapBlack
	return &ret
}

func (m *OrderedMap[K, V]) removeMax() (result, removed *OrderedMap[K, V]) {
	if m.right == nil {
		return m.remove(), m
	}
	right, removed := m.right.removeMax()
	return m.adopt(m.left, right).bubble(), removed
}

func (m *OrderedMap[K, V]) redden() *OrderedMap[K, V] {
	if m.color == orderedMapDoubleBlack && m.len == 0 {
		return nil
	}
	ret := *m
	ret.color--
	return &ret
}

func (m *OrderedMap[K, V]) bubble() *OrderedMap[K, V] {
	if (m.left != nil && m.left.color == orderedMapDoubleBlack) || (m.right != nil && m.right.color == orderedMapDoubleBlack) {
		unbalanced := &OrderedMap[K, V]{
			len:   m.len,
			color: m.color + 1,
			left:  m.left.redden(),
			right: m.right.redden(),
			key:   m.key,
			value: m.value,
		}
		if m.left != nil && m.left.color == orderedMapDoubleBlack {
			return unbalanced.balanceRight()
		}
		return unbalanced.balanceLeft()
	}
	return m
}

// OrderedMapElement represents a key-value pair and can be used to iterate over elements in a map.
type OrderedMapElement[K constraints.Ordered, V any] struct {
	lineage *Stack[*OrderedMap[K, V]]
	element *OrderedMap[K, V]
}

// Key returns the key of the represented element.
func (e *OrderedMapElement[K, V]) Key() K {
	return e.element.key
}

// Value returns the value of the represented element.
func (e *OrderedMapElement[K, V]) Value() V {
	return e.element.value
}

// Next returns the next element in the map.
//
// Complexity: O(log n) worst-case, amortized O(1) if iterating over the entire map
func (e *OrderedMapElement[K, V]) Next() *OrderedMapElement[K, V] {
	if !e.element.right.Empty() {
		lineage := e.lineage.Push(e.element)
		m := e.element.right
		for !m.Empty() && m.left != nil {
			lineage = lineage.Push(m)
			m = m.left
		}
		return &OrderedMapElement[K, V]{
			lineage: lineage,
			element: m,
		}
	}
	for l := e.lineage; !l.Empty(); l = l.Pop() {
		if e.element.key < l.Peek().key {
			return &OrderedMapElement[K, V]{
				lineage: l.Pop(),
				element: l.Peek(),
			}
		}
	}
	return nil
}

// Prev returns the previous element in the map.
//
// Complexity: O(log n) worst-case, amortized O(1) if iterating over an entire map
func (e *OrderedMapElement[K, V]) Prev() *OrderedMapElement[K, V] {
	if !e.element.left.Empty() {
		lineage := e.lineage.Push(e.element)
		m := e.element.left
		for !m.Empty() && m.right != nil {
			lineage = lineage.Push(m)
			m = m.right
		}
		return &OrderedMapElement[K, V]{
			lineage: lineage,
			element: m,
		}
	}
	for l := e.lineage; !l.Empty(); l = l.Pop() {
		if l.Peek().key < e.element.key {
			return &OrderedMapElement[K, V]{
				lineage: l.Pop(),
				element: l.Peek(),
			}
		}
	}
	return nil
}

// CountLess returns the number of elements that are less than this element.
//
// Complexity: O(log n) worst-case
func (e *OrderedMapElement[K, V]) CountLess() int {
	count := e.element.left.Len()
	for l := e.lineage; !l.Empty(); l = l.Pop() {
		if l.Peek().key < e.element.key {
			count += 1 + l.Peek().left.Len()
		}
	}
	return count
}

// CountGreater returns the number of elements that are greater than this element.
//
// Complexity: O(log n) worst-case
func (e *OrderedMapElement[K, V]) CountGreater() int {
	count := e.element.right.Len()
	for l := e.lineage; !l.Empty(); l = l.Pop() {
		if e.element.key < l.Peek().key {
			count += 1 + l.Peek().right.Len()
		}
	}
	return count
}
