package immutable

import (
	"strings"
)

const (
	orderedMapNegativeBlack = -1
	orderedMapRed           = 0
	orderedMapBlack         = 1
	orderedMapDoubleBlack   = 2
)

// OrderedMap implements an ordered map.
//
// Nil and the zero value for OrderedMap are both empty maps.
type OrderedMap struct {
	color        int
	left         *OrderedMap
	right        *OrderedMap
	key          interface{}
	value        interface{}
	lessThanFunc func(interface{}, interface{}) bool
}

// Empty returns true if the map is empty.
//
// Complexity: O(1) worst-case
func (m *OrderedMap) Empty() bool {
	return m == nil || m.key == nil
}

// Get returns the value associated with the given key if set.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Get(key interface{}) (interface{}, bool) {
	l := m.lessThanFuncOrEqual(key, nil)
	if l == nil {
		return nil, false
	}
	if !m.lessThanFunc(l.key, key) {
		return l.value, true
	}
	return nil, false
}

// Set associates a value with the given key.
//
// Only the built-in types may be used as keys. Once a value is set within a map, all subsequent
// operations must use the same key type.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Set(key, value interface{}) *OrderedMap {
	ret := m.insert(key, value)
	ret.color = orderedMapBlack
	return ret
}

// Delete removes a key from the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Delete(key interface{}) *OrderedMap {
	if ret, _ := m.delete(key); !ret.Empty() {
		ret.color = orderedMapBlack
		return ret
	}
	return nil
}

// Min returns the minimum element in the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Min() *OrderedMapElement {
	return m.min(nil)
}

// Max returns the maximum element in the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Max() *OrderedMapElement {
	return m.max(nil)
}

// MinAfter returns the minimum element in the map that is greater than the given key.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) MinAfter(key interface{}) *OrderedMapElement {
	return m.minGreaterThan(key, nil)
}

// MaxBefore returns the maximum element in the map that is less than the given key.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) MaxBefore(key interface{}) *OrderedMapElement {
	return m.maxLessThan(key, nil)
}

func (m *OrderedMap) min(lineage *Stack) *OrderedMapElement {
	if m.Empty() {
		return nil
	} else if m.left != nil {
		return m.left.min(lineage.Push(m))
	}
	return &OrderedMapElement{
		lineage: lineage,
		element: m,
	}
}

func (m *OrderedMap) max(lineage *Stack) *OrderedMapElement {
	if m.Empty() {
		return nil
	} else if m.right != nil {
		return m.right.max(lineage.Push(m))
	}
	return &OrderedMapElement{
		lineage: lineage,
		element: m,
	}
}

func (m *OrderedMap) minGreaterThan(key interface{}, lineage *Stack) *OrderedMapElement {
	if m.Empty() {
		return nil
	} else if m.lessThanFunc(key, m.key) {
		if m.left != nil {
			if r := m.left.minGreaterThan(key, lineage.Push(m)); r != nil {
				return r
			}
		}
		return &OrderedMapElement{
			lineage: lineage,
			element: m,
		}
	} else if m.lessThanFunc(m.key, key) {
		return m.right.minGreaterThan(key, lineage.Push(m))
	}
	return m.right.min(lineage.Push(m))
}

func (m *OrderedMap) maxLessThan(key interface{}, lineage *Stack) *OrderedMapElement {
	if m.Empty() {
		return nil
	} else if m.lessThanFunc(m.key, key) {
		if m.right != nil {
			if r := m.right.maxLessThan(key, lineage.Push(m)); r != nil {
				return r
			}
		}
		return &OrderedMapElement{
			lineage: lineage,
			element: m,
		}
	} else if m.lessThanFunc(key, m.key) {
		return m.left.maxLessThan(key, lineage.Push(m))
	}
	return m.left.max(lineage.Push(m))
}

func (m *OrderedMap) delete(key interface{}) (*OrderedMap, bool) {
	if m.Empty() {
		return m, false
	} else if m.lessThanFunc(key, m.key) {
		if left, didDelete := m.left.delete(key); didDelete {
			return m.adopt(left, m.right).bubble(), true
		}
		return m, false
	} else if m.lessThanFunc(m.key, key) {
		if right, didDelete := m.right.delete(key); didDelete {
			return m.adopt(m.left, right).bubble(), true
		}
		return m, false
	}
	return m.remove(), true
}

func (m *OrderedMap) adopt(left, right *OrderedMap) *OrderedMap {
	return &OrderedMap{
		color:        m.color,
		left:         left,
		right:        right,
		key:          m.key,
		value:        m.value,
		lessThanFunc: m.lessThanFunc,
	}
}

func (m *OrderedMap) lessThanFuncOrEqual(key interface{}, candidate *OrderedMap) *OrderedMap {
	if m.Empty() {
		return candidate
	} else if m.lessThanFunc(key, m.key) {
		return m.left.lessThanFuncOrEqual(key, candidate)
	}
	return m.right.lessThanFuncOrEqual(key, m)
}

func (m *OrderedMap) insert(key, value interface{}) *OrderedMap {
	if m.Empty() {
		return &OrderedMap{
			color:        orderedMapRed,
			key:          key,
			value:        value,
			lessThanFunc: builtInLessThan(key),
		}
	} else if m.lessThanFunc(key, m.key) {
		return m.adopt(m.left.insert(key, value), m.right).balanceLeft()
	} else if m.lessThanFunc(m.key, key) {
		return m.adopt(m.left, m.right.insert(key, value)).balanceRight()
	}
	return &OrderedMap{
		color:        m.color,
		left:         m.left,
		right:        m.right,
		key:          m.key,
		value:        value,
		lessThanFunc: m.lessThanFunc,
	}
}

func (m *OrderedMap) balanceLeft() *OrderedMap {
	if m.color >= orderedMapBlack && m.left != nil {
		if m.left.color == orderedMapRed {
			if m.left.left != nil && m.left.left.color == orderedMapRed {
				return &OrderedMap{
					color: m.color - 1,
					left: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.left.left.left,
						right:        m.left.left.right,
						key:          m.left.left.key,
						value:        m.left.left.value,
						lessThanFunc: m.lessThanFunc,
					},
					right: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.left.right,
						right:        m.right,
						key:          m.key,
						value:        m.value,
						lessThanFunc: m.lessThanFunc,
					},
					key:          m.left.key,
					value:        m.left.value,
					lessThanFunc: m.lessThanFunc,
				}
			} else if m.left.right != nil && m.left.right.color == orderedMapRed {
				return &OrderedMap{
					color: m.color - 1,
					left: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.left.left,
						right:        m.left.right.left,
						key:          m.left.key,
						value:        m.left.value,
						lessThanFunc: m.lessThanFunc,
					},
					right: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.left.right.right,
						right:        m.right,
						key:          m.key,
						value:        m.value,
						lessThanFunc: m.lessThanFunc,
					},
					key:          m.left.right.key,
					value:        m.left.right.value,
					lessThanFunc: m.lessThanFunc,
				}
			}
		} else if m.left.color == orderedMapNegativeBlack {
			unbalancedLeft := &OrderedMap{
				color:        orderedMapBlack,
				left:         m.left.left.redden(),
				right:        m.left.right.left,
				key:          m.left.key,
				value:        m.left.value,
				lessThanFunc: m.lessThanFunc,
			}
			return &OrderedMap{
				color: orderedMapBlack,
				left:  unbalancedLeft.balanceLeft(),
				right: &OrderedMap{
					color:        orderedMapBlack,
					left:         m.left.right.right,
					right:        m.right,
					key:          m.key,
					value:        m.value,
					lessThanFunc: m.lessThanFunc,
				},
				key:          m.left.right.key,
				value:        m.left.right.value,
				lessThanFunc: m.lessThanFunc,
			}
		}
	}
	return m
}

func (m *OrderedMap) balanceRight() *OrderedMap {
	if m.color >= orderedMapBlack && m.right != nil {
		if m.right.color == orderedMapRed {
			if m.right.left != nil && m.right.left.color == orderedMapRed {
				return &OrderedMap{
					color: m.color - 1,
					left: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.left,
						right:        m.right.left.left,
						key:          m.key,
						value:        m.value,
						lessThanFunc: m.lessThanFunc,
					},
					right: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.right.left.right,
						right:        m.right.right,
						key:          m.right.key,
						value:        m.right.value,
						lessThanFunc: m.lessThanFunc,
					},
					key:          m.right.left.key,
					value:        m.right.left.value,
					lessThanFunc: m.lessThanFunc,
				}
			} else if m.right.right != nil && m.right.right.color == orderedMapRed {
				return &OrderedMap{
					color: m.color - 1,
					left: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.left,
						right:        m.right.left,
						key:          m.key,
						value:        m.value,
						lessThanFunc: m.lessThanFunc,
					},
					right: &OrderedMap{
						color:        orderedMapBlack,
						left:         m.right.right.left,
						right:        m.right.right.right,
						key:          m.right.right.key,
						value:        m.right.right.value,
						lessThanFunc: m.lessThanFunc,
					},
					key:          m.right.key,
					value:        m.right.value,
					lessThanFunc: m.lessThanFunc,
				}
			}
		} else if m.right.color == orderedMapNegativeBlack {
			unbalancedRight := &OrderedMap{
				color:        orderedMapBlack,
				left:         m.right.left.right,
				right:        m.right.right.redden(),
				key:          m.right.key,
				value:        m.right.value,
				lessThanFunc: m.lessThanFunc,
			}
			return &OrderedMap{
				color: orderedMapBlack,
				left: &OrderedMap{
					color:        orderedMapBlack,
					left:         m.left,
					right:        m.right.left.left,
					key:          m.key,
					value:        m.value,
					lessThanFunc: m.lessThanFunc,
				},
				right:        unbalancedRight.balanceRight(),
				key:          m.right.left.key,
				value:        m.right.left.value,
				lessThanFunc: m.lessThanFunc,
			}
		}
	}
	return m
}

var doubleBlackLeaf = &OrderedMap{color: orderedMapDoubleBlack}

func (m *OrderedMap) remove() *OrderedMap {
	if !m.left.Empty() && !m.right.Empty() {
		left, removed := m.left.removeMax()
		reduced := &OrderedMap{
			color:        m.color,
			left:         left,
			right:        m.right,
			key:          removed.key,
			value:        removed.value,
			lessThanFunc: m.lessThanFunc,
		}
		return reduced.bubble()
	}
	var child *OrderedMap
	if !m.left.Empty() {
		child = m.left
	} else if !m.right.Empty() {
		child = m.right
	} else {
		if m.color == orderedMapRed {
			return nil
		}
		return doubleBlackLeaf
	}
	ret := *child
	ret.color = orderedMapBlack
	return &ret
}

func (m *OrderedMap) removeMax() (result, removed *OrderedMap) {
	if m.right == nil {
		return m.remove(), m
	}
	right, removed := m.right.removeMax()
	return m.adopt(m.left, right).bubble(), removed
}

func (m *OrderedMap) redden() *OrderedMap {
	if m == doubleBlackLeaf {
		return nil
	}
	ret := *m
	ret.color--
	return &ret
}

func (m *OrderedMap) bubble() *OrderedMap {
	if (m.left != nil && m.left.color == orderedMapDoubleBlack) || (m.right != nil && m.right.color == orderedMapDoubleBlack) {
		unbalanced := &OrderedMap{
			color:        m.color + 1,
			left:         m.left.redden(),
			right:        m.right.redden(),
			key:          m.key,
			value:        m.value,
			lessThanFunc: m.lessThanFunc,
		}
		if m.left != nil && m.left.color == orderedMapDoubleBlack {
			return unbalanced.balanceRight()
		}
		return unbalanced.balanceLeft()
	}
	return m
}

// OrderedMapElement represents a key-value pair and can be used to iterate over elements in a map.
type OrderedMapElement struct {
	lineage *Stack
	element *OrderedMap
}

// Key returns the key of the represented element.
func (e *OrderedMapElement) Key() interface{} {
	return e.element.key
}

// Value returns the value of the represented element.
func (e *OrderedMapElement) Value() interface{} {
	return e.element.value
}

// Next returns the next element in the map.
//
// Complexity: O(log n) worst-case, amortized O(1) if iterating over the entire map
func (e *OrderedMapElement) Next() *OrderedMapElement {
	if !e.element.right.Empty() {
		lineage := e.lineage.Push(e.element)
		m := e.element.right
		for !m.Empty() && m.left != nil {
			lineage = lineage.Push(m)
			m = m.left
		}
		return &OrderedMapElement{
			lineage: lineage,
			element: m,
		}
	}
	for l := e.lineage; !l.Empty(); l = l.Pop() {
		if e.element.lessThanFunc(e.element.key, l.Peek().(*OrderedMap).key) {
			return &OrderedMapElement{
				lineage: l.Pop(),
				element: l.Peek().(*OrderedMap),
			}
		}
	}
	return nil
}

// Prev returns the previous element in the map.
//
// Complexity: O(log n) worst-case, amortized O(1) if iterating over an entire map
func (e *OrderedMapElement) Prev() *OrderedMapElement {
	if !e.element.left.Empty() {
		lineage := e.lineage.Push(e.element)
		m := e.element.left
		for !m.Empty() && m.right != nil {
			lineage = lineage.Push(m)
			m = m.right
		}
		return &OrderedMapElement{
			lineage: lineage,
			element: m,
		}
	}
	for l := e.lineage; !l.Empty(); l = l.Pop() {
		if e.element.lessThanFunc(l.Peek().(*OrderedMap).key, e.element.key) {
			return &OrderedMapElement{
				lineage: l.Pop(),
				element: l.Peek().(*OrderedMap),
			}
		}
	}
	return nil
}

func builtInLessThan(value interface{}) func(interface{}, interface{}) bool {
	switch value.(type) {
	case int:
		return func(a, b interface{}) bool { return a.(int) < b.(int) }
	case int8:
		return func(a, b interface{}) bool { return a.(int8) < b.(int8) }
	case int16:
		return func(a, b interface{}) bool { return a.(int16) < b.(int16) }
	case int32:
		return func(a, b interface{}) bool { return a.(int32) < b.(int32) }
	case int64:
		return func(a, b interface{}) bool { return a.(int64) < b.(int64) }
	case uint:
		return func(a, b interface{}) bool { return a.(uint) < b.(uint) }
	case uint8:
		return func(a, b interface{}) bool { return a.(uint8) < b.(uint8) }
	case uint16:
		return func(a, b interface{}) bool { return a.(uint16) < b.(uint16) }
	case uint32:
		return func(a, b interface{}) bool { return a.(uint32) < b.(uint32) }
	case uint64:
		return func(a, b interface{}) bool { return a.(uint64) < b.(uint64) }
	case uintptr:
		return func(a, b interface{}) bool { return a.(uintptr) < b.(uintptr) }
	case float32:
		return func(a, b interface{}) bool { return a.(float32) < b.(float32) }
	case float64:
		return func(a, b interface{}) bool { return a.(float64) < b.(float64) }
	case string:
		return func(a, b interface{}) bool { return strings.Compare(a.(string), b.(string)) == -1 }
	}
	panic("invalid type")
}
