package data

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type Set[K fmt.Stringer] map[string]K

func NewSet[K fmt.Stringer](ms ...map[string]K) Set[K] {
	dst := map[string]K{}
	for _, src := range ms {
		maps.Copy(dst, src)
	}

	return dst
}

func (s Set[K]) Get(key K) (K, bool) {
	v, ok := s[key.String()]
	return v, ok
}

func (s Set[K]) GetString(str string) (K, bool) {
	v, ok := s[str]
	return v, ok
}

func (s Set[K]) Set(key K) bool {
	_, exists := s[key.String()]
	s[key.String()] = key
	return !exists
}

func (s Set[K]) Delete(key K) {
	delete(s, key.String())
}

func (s Set[K]) Values() []K {
	return maps.Values(s)
}

func (s Set[K]) Keys() []string {
	return maps.Keys(s)
}

func (m Set[K]) Map() map[string]K {
	return m
}
