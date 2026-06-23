package service

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"gorm.io/gorm"
)

type IPManager struct {
	DB         *gorm.DB
	pools      map[string]*IPPool
	mu         sync.RWMutex
}

type IPPool struct {
	Network    net.IPNet
	UsedIPs    map[string]bool
	LastUsed   int
}

func NewIPManager(db *gorm.DB) *IPManager {
	return &IPManager{
		DB:  db,
		pools: make(map[string]*IPPool),
	}
}

func (m *IPManager) AllocateIP(hostID uint) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	poolName := fmt.Sprintf("host-%d", hostID)
	
	pool, exists := m.pools[poolName]
	if !exists {
		// Create new pool for host
		_, network, err := net.ParseCIDR("10.0.0.0/24")
		if err != nil {
			return "", err
		}
		
		pool = &IPPool{
			Network: *network,
			UsedIPs: make(map[string]bool),
		}
		m.pools[poolName] = pool
	}

	// Find available IP
	start := pool.LastUsed + 1
	for i := 0; i < 254; i++ {
		idx := (start + i) % 254 + 1
		ip := net.IPv4(10, 0, byte(hostID), byte(idx))
		ipStr := ip.String()
		
		if !pool.UsedIPs[ipStr] {
			pool.UsedIPs[ipStr] = true
			pool.LastUsed = idx
			return ipStr, nil
		}
	}

	return "", errors.New("no available IPs in pool")
}

func (m *IPManager) ReleaseIP(hostID uint, ip string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	poolName := fmt.Sprintf("host-%d", hostID)
	if pool, exists := m.pools[poolName]; exists {
		delete(pool.UsedIPs, ip)
	}
	return nil
}

func (m *IPManager) GetAvailableIPs(hostID uint) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	poolName := fmt.Sprintf("host-%d", hostID)
	pool, exists := m.pools[poolName]
	if !exists {
		return []string{}, nil
	}

	available := make([]string, 0)
	for i := 1; i <= 254; i++ {
		ip := net.IPv4(10, 0, byte(hostID), byte(i))
		ipStr := ip.String()
		if !pool.UsedIPs[ipStr] {
			available = append(available, ipStr)
		}
	}

	return available, nil
}
