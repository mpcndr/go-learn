//go:build solution

package generics

func Map[T, U any](s []T, fn func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, fn(v))
	}
	return out
}

func Filter[T any](s []T, keep func(T) bool) []T {
	var out []T
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func Reduce[T, A any](s []T, initial A, fn func(A, T) A) A {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func GroupBy[T any, K comparable](s []T, keyFn func(T) K) map[K][]T {
	out := make(map[K][]T)
	for _, v := range s {
		k := keyFn(v)
		out[k] = append(out[k], v)
	}
	return out
}

func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(items))}
	for _, v := range items {
		s.Add(v)
	}
	return s
}

func (s *Set[T]) Add(v T)      { s.m[v] = struct{}{} }
func (s *Set[T]) Has(v T) bool { _, ok := s.m[v]; return ok }
func (s *Set[T]) Len() int     { return len(s.m) }

func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	out := NewSet[T]()
	for v := range s.m {
		if other.Has(v) {
			out.Add(v)
		}
	}
	return out
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{capacity: capacity, data: make(map[K]V, capacity)}
}

func (c *Cache[K, V]) Put(key K, value V) {
	if _, exists := c.data[key]; exists {
		c.data[key] = value
		return
	}
	if len(c.order) >= c.capacity {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.data, oldest)
	}
	c.data[key] = value
	c.order = append(c.order, key)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	v, ok := c.data[key]
	return v, ok
}
