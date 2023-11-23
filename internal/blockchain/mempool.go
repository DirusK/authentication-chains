/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package blockchain

import (
	"bytes"
	"sync"

	"authentication-chains/internal/types"
)

//go:generate ifacemaker -f mempool.go -s memPool -p blockchain -i MemPool -y "MemPool - describe an interface for working with memory pool."

// memPool implements mem-pool logic.
type memPool struct {
	mutex   sync.RWMutex
	memPool []*types.DeviceAuthenticationRequest
}

// NewMemPool creates a new mem-pool instance.
func newMemPool() MemPool {
	return &memPool{
		memPool: make([]*types.DeviceAuthenticationRequest, 0),
	}
}

// GetFirst returns the first device authentication request from the mem-pool.
func (m *memPool) GetFirst() *types.DeviceAuthenticationRequest {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.memPool) == 0 {
		return nil
	}

	return m.memPool[0]
}

// GetAll returns all device authentication requests from the mem-pool.
func (m *memPool) GetAll() []*types.DeviceAuthenticationRequest {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]*types.DeviceAuthenticationRequest, len(m.memPool))
	copy(result, m.memPool)

	return result
}

// Add adds a device authentication request to the mem-pool.
func (m *memPool) Add(request *types.DeviceAuthenticationRequest) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.memPool = append(m.memPool, request)
}

// Remove removes a device authentication request from the mem-pool.
func (m *memPool) Remove(request *types.DeviceAuthenticationRequest) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, r := range m.memPool {
		if bytes.Equal(r.DeviceId, request.DeviceId) {
			m.memPool = append(m.memPool[:i], m.memPool[i+1:]...)
			break
		}
	}
}
