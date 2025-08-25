package pipeline

import (
	"sync"
	"github.com/jstaud/edgelink-go/pkg/models"
)

// MemoryCache stores the latest reading from each device in memory
type MemoryCache struct {
	mu sync.RWMutex                   // Read-Write mutex for thread safety
	m  map[string]models.Reading      // Map from device ID to latest reading
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		m: make(map[string]models.Reading),  // Initialize empty map
	}
}

// Set stores a reading (thread-safe)
func (c *MemoryCache) Set(r models.Reading) {
	c.mu.Lock()           // Get exclusive write lock
	c.m[r.DeviceID] = r   // Store the reading
	c.mu.Unlock()         // Release the lock
}

// Latest gets the most recent reading for a device (thread-safe)
func (c *MemoryCache) Latest(id string) (models.Reading, bool) {
	c.mu.RLock()          // Get shared read lock (multiple readers OK)
	defer c.mu.RUnlock()  // Release lock when function exits
	
	reading, exists := c.m[id]  // Look up reading, get value + whether it exists
	return reading, exists
}
