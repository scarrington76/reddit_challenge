package services

import (
	"maps"
	"sync"
)

type PostStats struct {
	User  string
	Ups   int
	Title string
}

type SafeMap struct {
	sync.Mutex
	data map[string]PostStats
}

// NewSafeMap creates a new safe map.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]PostStats),
	}
}

// Get sets a value in the map safely.
func (m *SafeMap) Get(key string) (PostStats, bool) {
	m.Lock()
	defer m.Unlock()
	val, ok := m.data[key]
	return val, ok
}

// Set sets a value in the map safely.
func (m *SafeMap) Set(key string, value PostStats) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = value
}

// Length calculates and returns the length of the underlying map safely.
func (m *SafeMap) Length() int {
	m.Lock()
	defer m.Unlock()
	return len(m.data)
}

// ClonePostMap returns a shallow clone of the map safely. This is intended
// for endpoints which need to access data.
func (m *SafeMap) ClonePostMap() map[string]PostStats {
	m.Lock()
	defer m.Unlock()
	return maps.Clone(m.data)
}
