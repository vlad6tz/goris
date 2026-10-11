package store

import (
	"bytes"
	"sync"
)

type Store struct {
	mu sync.RWMutex
	m  map[string][]byte
}

func New() *Store {
	return &Store{
		m: make(map[string][]byte),
	}
}
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock() // read locking around a map read
	v, ok := s.m[key]
	s.mu.RUnlock()
	return v, ok
}
func (s *Store) Set(key string, value []byte) {
	s.mu.Lock() // write lock around a map write
	s.m[key] = bytes.Clone(value)
	s.mu.Unlock()
}
