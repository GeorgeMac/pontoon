package monitor

import (
	"fmt"
	"sync"
)

type Monitor struct {
	mo map[string]Trackable
	mu *sync.RWMutex
}

func NewMonitor() *Monitor {
	return &Monitor{
		mo: make(map[string]Trackable),
		mu: &sync.RWMutex{},
	}
}

func (m *Monitor) Get(id string) (t Trackable, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var ok bool
	if t, ok = m.mo[id]; !ok {
		err = TrackableNotExistError{id}
		return
	}
	return
}

func (m *Monitor) Put(id string, t Trackable) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.mo[id]; ok {
		err = TrackableExistsError{id}
		return
	}

	m.mo[id] = t
	return
}

func (m *Monitor) Status(id string) Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.mo[id]
	if !ok {
		return UNKNOWN
	}

	return t.Status()
}

type TrackableExistsError struct {
	id string
}

func (t TrackableExistsError) Error() string {
	return fmt.Sprintf("Trackable with id %s already exists", t.id)
}

type TrackableNotExistError struct {
	id string
}

func (t TrackableNotExistError) Error() string {
	return fmt.Sprintf("Trackable with id %s does not exists", t.id)
}