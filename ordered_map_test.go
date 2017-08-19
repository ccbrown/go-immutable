package immutable

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap(t *testing.T) {
	var m *OrderedMap
	assert.True(t, m.Empty())

	m = m.Set("foo", "bar")
	assert.False(t, m.Empty())
	v, ok := m.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", v)

	_, ok = m.Get("fom")
	assert.False(t, ok)

	_, ok = m.Get("fop")
	assert.False(t, ok)

	m = m.Set("qux", "quux")
	assert.False(t, m.Empty())
	v, ok = m.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", v)
	v, ok = m.Get("qux")
	assert.True(t, ok)
	assert.Equal(t, "quux", v)
}

func TestOrderedMap_Iteration(t *testing.T) {
	var m *OrderedMap
	for i := 0; i < 1000; i++ {
		m = m.Set(i, i*2)
	}

	e := m.Min()
	for i := 0; i < 1000; i++ {
		assert.NotNil(t, e)
		assert.Equal(t, i, e.Key())
		assert.Equal(t, i*2, e.Value())
		e = e.Next()
	}
	assert.Nil(t, e)

	e = m.Max()
	for i := 999; i >= 0; i-- {
		assert.NotNil(t, e)
		assert.Equal(t, i, e.Key())
		assert.Equal(t, i*2, e.Value())
		e = e.Prev()
	}
	assert.Nil(t, e)
}

func TestOrderedMap_Fuzz(t *testing.T) {
	ref := make(map[int]int)
	var m *OrderedMap
	assert.True(t, m.Empty())
	for i := 0; i < 10000; i++ {
		k := rand.Int()
		v := rand.Int()
		ref[k] = v
		m = m.Set(k, v)
	}
	for k, refv := range ref {
		v, ok := m.Get(k)
		assert.True(t, ok)
		assert.Equal(t, refv, v)
	}
}

var orderedMapValueResult interface{}

func BenchmarkOrderedMap_Get(b *testing.B) {
	for _, n := range []int{100, 10000, 1000000} {
		m := &OrderedMap{}
		for i := 0; i < n; i++ {
			m = m.Set(i, "foo")
		}
		b.Run(fmt.Sprintf("n=%v", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, _ := m.Get(i % n)
				orderedMapValueResult = v
			}
		})
	}
}

var orderedMapResult *OrderedMap

func BenchmarkOrderedMap_Set(b *testing.B) {
	for _, n := range []int{100, 10000, 1000000} {
		m := &OrderedMap{}
		for i := 0; i < n; i++ {
			m = m.Set(i, "foo")
		}
		b.Run(fmt.Sprintf("n=%v", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				orderedMapResult = m.Set(i%n, "bar")
			}
		})
	}
}
