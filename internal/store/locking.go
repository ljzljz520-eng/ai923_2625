package store

import (
	"errors"
	"sync"
	"time"
)

type Lease struct {
	ID        string
	ExpiresAt time.Time
}
type LeaseManager struct {
	mu     sync.Mutex
	leases map[string]Lease
	now    func() time.Time
}

func NewLeaseManager() *LeaseManager { return &LeaseManager{leases: map[string]Lease{}, now: time.Now} }
func (m *LeaseManager) Acquire(key string, duration time.Duration) (Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if key == "" {
		return Lease{}, errors.New("lease key is required")
	}
	if current, ok := m.leases[key]; ok && current.ExpiresAt.After(m.now()) {
		return Lease{}, errors.New("lease is already held")
	}
	lease := Lease{ID: key, ExpiresAt: m.now().Add(duration)}
	m.leases[key] = lease
	return lease, nil
}
func (m *LeaseManager) Release(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.leases[key]; !ok {
		return false
	}
	delete(m.leases, key)
	return true
}
func (m *LeaseManager) Held(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	lease, ok := m.leases[key]
	return ok && lease.ExpiresAt.After(m.now())
}
func (m *LeaseManager) Expire() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	removed := 0
	now := m.now()
	for key, lease := range m.leases {
		if !lease.ExpiresAt.After(now) {
			delete(m.leases, key)
			removed++
		}
	}
	return removed
}
func (m *LeaseManager) Count() int { m.mu.Lock(); defer m.mu.Unlock(); return len(m.leases) }
