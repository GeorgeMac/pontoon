package jobs

import (
	"fmt"
	"github.com/GeorgeMac/pontoon/monitor"
	"sync"
)

type Store struct {
	mo map[string]*Job
	mu *sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		mo: make(map[string]*Job),
		mu: &sync.RWMutex{},
	}
}

func (m *Store) Get(id string) (t *Job, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var ok bool
	if t, ok = m.mo[id]; !ok {
		err = DoesNotExistError{id}
		return
	}
	return
}

func (m *Store) Put(id string, t *Job) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.mo[id]; ok {
		err = ExistsError{id}
		return
	}

	m.mo[id] = t
	return
}

func (m *Store) Report(id string) monitor.Report {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.mo[id]
	if !ok {
		return monitor.Report{
			Name:   id,
			Status: monitor.UNKNOWN.String(),
		}
	}

	return t.Report()
}

func (m *Store) FullReport(id string) monitor.FullReport {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.mo[id]
	if !ok {
		return monitor.FullReport{
			Report: monitor.Report{
				Name:   id,
				Status: monitor.UNKNOWN.String(),
			},
		}
	}

	return monitor.FullReport{
		Report:  t.Report(),
		History: t.History(),
	}
}

func (m *Store) List() (t []monitor.Report) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t = make([]monitor.Report, 0)

	for _, v := range m.mo {
		t = append(t, v.Report())
	}

	return
}

type ExistsError struct {
	id string
}

func (t ExistsError) Error() string {
	return fmt.Sprintf("Entry with id %s already exists", t.id)
}

type DoesNotExistError struct {
	id string
}

func (t DoesNotExistError) Error() string {
	return fmt.Sprintf("Entry with id %s does not exists", t.id)
}
