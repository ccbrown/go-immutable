package immutable

import (
	"strings"
)

// OrderedMap implements an ordered map.
//
// Nil and the zero value for OrderedMap are both empty maps.
type OrderedMap struct {
	red      bool
	left     *OrderedMap
	right    *OrderedMap
	key      interface{}
	value    interface{}
	lessThan func(interface{}, interface{}) bool
}

// Returns true if the map is empty.
//
// Complexity: O(1) worst-case
func (m *OrderedMap) Empty() bool {
	return m == nil || m.key == nil
}

// Returns the value associated with the given key if set.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Get(key interface{}) (interface{}, bool) {
	l := m.lessThanOrEqual(key, nil)
	if l == nil {
		return nil, false
	}
	if !m.lessThan(l.key, key) {
		return l.value, true
	}
	return nil, false
}

// Associates a value with the given key.
//
// Only the built-in types may be used as keys. Once a value is set within a map, all subsequent
// operations must use the same key type.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Set(key, value interface{}) *OrderedMap {
	ret := m.insert(key, value)
	ret.red = false
	return ret
}

// Returns the smallest element in the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Min() *OrderedMapElement {
	var lineage *Stack
	for !m.Empty() && m.left != nil {
		lineage = lineage.Push(m)
		m = m.left
	}
	return &OrderedMapElement{
		lineage: lineage,
		element: m,
	}
}

// Returns the largest element in the map.
//
// Complexity: O(log n) worst-case
func (m *OrderedMap) Max() *OrderedMapElement {
	var lineage *Stack
	for !m.Empty() && m.right != nil {
		lineage = lineage.Push(m)
		m = m.right
	}
	return &OrderedMapElement{
		lineage: lineage,
		element: m,
	}
}

func (m *OrderedMap) lessThanOrEqual(key interface{}, candidate *OrderedMap) *OrderedMap {
	if m.Empty() {
		return candidate
	} else if m.lessThan(key, m.key) {
		return m.left.lessThanOrEqual(key, candidate)
	}
	return m.right.lessThanOrEqual(key, m)
}

func (m *OrderedMap) insert(key, value interface{}) *OrderedMap {
	if m.Empty() {
		return &OrderedMap{
			red:      true,
			key:      key,
			value:    value,
			lessThan: buildInLessThan(key),
		}
	} else if m.lessThan(key, m.key) {
		return m.balanceLeft(m.left.insert(key, value), m.right)
	} else if m.lessThan(m.key, key) {
		return m.balanceRight(m.left, m.right.insert(key, value))
	}
	return &OrderedMap{
		red:      m.red,
		left:     m.left,
		right:    m.right,
		key:      m.key,
		value:    value,
		lessThan: m.lessThan,
	}
}

func (m *OrderedMap) balanceLeft(left, right *OrderedMap) *OrderedMap {
	if !m.red && left != nil && left.red {
		if left.left != nil && left.left.red {
			return &OrderedMap{
				red: true,
				left: &OrderedMap{
					red:      false,
					left:     left.left.left,
					right:    left.left.right,
					key:      left.left.key,
					value:    left.left.value,
					lessThan: m.lessThan,
				},
				right: &OrderedMap{
					red:      false,
					left:     left.right,
					right:    right,
					key:      m.key,
					value:    m.value,
					lessThan: m.lessThan,
				},
				key:      left.key,
				value:    left.value,
				lessThan: m.lessThan,
			}
		} else if left.right != nil && left.right.red {
			return &OrderedMap{
				red: true,
				left: &OrderedMap{
					red:      false,
					left:     left.left,
					right:    left.right.left,
					key:      left.key,
					value:    left.value,
					lessThan: m.lessThan,
				},
				right: &OrderedMap{
					red:      false,
					left:     left.right.right,
					right:    right,
					key:      m.key,
					value:    m.value,
					lessThan: m.lessThan,
				},
				key:      left.right.key,
				value:    left.right.value,
				lessThan: m.lessThan,
			}
		}
	}
	return &OrderedMap{
		red:      m.red,
		left:     left,
		right:    right,
		key:      m.key,
		value:    m.value,
		lessThan: m.lessThan,
	}
}

func (m *OrderedMap) balanceRight(left, right *OrderedMap) *OrderedMap {
	if !m.red && right != nil && right.red {
		if right.left != nil && right.left.red {
			return &OrderedMap{
				red: true,
				left: &OrderedMap{
					red:      false,
					left:     left,
					right:    right.left.left,
					key:      m.key,
					value:    m.value,
					lessThan: m.lessThan,
				},
				right: &OrderedMap{
					red:      false,
					left:     right.left.right,
					right:    right.right,
					key:      right.key,
					value:    right.value,
					lessThan: m.lessThan,
				},
				key:      right.left.key,
				value:    right.left.value,
				lessThan: m.lessThan,
			}
		} else if right.right != nil && right.right.red {
			return &OrderedMap{
				red: true,
				left: &OrderedMap{
					red:      false,
					left:     left,
					right:    right.left,
					key:      m.key,
					value:    m.value,
					lessThan: m.lessThan,
				},
				right: &OrderedMap{
					red:      false,
					left:     right.right.left,
					right:    right.right.right,
					key:      right.right.key,
					value:    right.right.value,
					lessThan: m.lessThan,
				},
				key:      right.key,
				value:    right.value,
				lessThan: m.lessThan,
			}
		}
	}
	return &OrderedMap{
		red:      m.red,
		left:     left,
		right:    right,
		key:      m.key,
		value:    m.value,
		lessThan: m.lessThan,
	}
}

// OrderedMapElement represents a key-value pair and can be used to iterate over elements in a map.
type OrderedMapElement struct {
	lineage *Stack
	element *OrderedMap
}

// Returns the key of the represented element.
func (e *OrderedMapElement) Key() interface{} {
	return e.element.key
}

// Returns the value of the represented element.
func (e *OrderedMapElement) Value() interface{} {
	return e.element.value
}

// Returns the next element in the map.
//
// Complexity: O(log n) worst-case, amortized O(1) if iterating over an entire map
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
		if e.element.lessThan(e.element.key, l.Peek().(*OrderedMap).key) {
			return &OrderedMapElement{
				lineage: l.Pop(),
				element: l.Peek().(*OrderedMap),
			}
		}
	}
	return nil
}

// Returns the previous element in the map.
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
		if e.element.lessThan(l.Peek().(*OrderedMap).key, e.element.key) {
			return &OrderedMapElement{
				lineage: l.Pop(),
				element: l.Peek().(*OrderedMap),
			}
		}
	}
	return nil
}

func buildInLessThan(value interface{}) func(interface{}, interface{}) bool {
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
