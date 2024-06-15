package data

import (
	"fmt"
	"sync"

	"golang.org/x/exp/maps"
)

type AsyncMap[K fmt.Stringer, V any] struct {
	m  map[string]V
	mu sync.RWMutex
}

func NewAsyncMap[K fmt.Stringer, V any](ms ...map[string]V) *AsyncMap[K, V] {
	dst := map[string]V{}
	for _, src := range ms {
		maps.Copy(dst, src)
	}
	return &AsyncMap[K, V]{m: dst}
}

func (m *AsyncMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.m[key.String()]
	return v, ok
}

func (m *AsyncMap[K, V]) Set(key K, value V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.m[key.String()]
	m.m[key.String()] = value
	return !exists
}

func (m *AsyncMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key.String())
}

func (m *AsyncMap[K, V]) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return maps.Keys(m.m)
}

func (m *AsyncMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return maps.Values(m.m)
}

func (m *AsyncMap[K, V]) Map() map[string]V {
	return m.m
}

type AsyncSet[K fmt.Stringer] AsyncMap[K, K]

func NewAsyncSet[K fmt.Stringer](ms ...map[string]K) *AsyncSet[K] {
	dst := map[string]K{}
	for _, src := range ms {
		maps.Copy(dst, src)
	}

	return &AsyncSet[K]{m: dst}
}

func (s *AsyncSet[K]) Get(key K) (K, bool) {
	return (*AsyncMap[K, K])(s).Get(key)
}
func (s *AsyncSet[K]) GetString(str string) (K, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[str]
	return v, ok
}

func (s *AsyncSet[K]) Set(key K) bool {
	return (*AsyncMap[K, K])(s).Set(key, key)
}

func (s *AsyncSet[K]) Delete(key K) {
	(*AsyncMap[K, K])(s).Delete(key)
}

func (s *AsyncSet[K]) Values() []K {
	return (*AsyncMap[K, K])(s).Values()
}

func (s *AsyncSet[K]) Keys() []string {
	return (*AsyncMap[K, K])(s).Keys()
}

func (m *AsyncSet[K]) Map() map[string]K {
	return m.m
}
