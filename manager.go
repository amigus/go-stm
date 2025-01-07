package stm

import (
	"sync"
	"time"
)

// TokenManager issues tokens and checks their validity.
type TokenManager interface {
	// Get returns a valid token subject to the maxUses and maxTime constraints.
	Get() string
	Check(token string) bool
}

// manager struct represents a "ring buffer."
// It holds the list of tokens, parameters, an index and a mutux.
type manager struct {
	tokens   []token
	maxCount int
	maxUses  int
	maxTime  time.Duration
	index    int
	mu       sync.Mutex
}

// Get returns a valid token subject to the maxUses and maxTime constraints.
func (m *manager) Get() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	token := &m.tokens[m.index]
	token.count++
	if m.maxTime > 0 && time.Now().After(token.expires) || token.count > m.maxUses {
		token.count = 1
		token.expires = time.Now().Add(m.maxTime)
		token.value = UUIDToken()
	}
	m.index = (m.index + 1) % m.maxCount
	return token.value
}

// Check returns true if the token is valid after incrementing the counter on it.
func (m *manager) Check(token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, t := range m.tokens {
		if t.value == token {
			m.tokens[i].count++
			return true
		}
	}
	return false
}
