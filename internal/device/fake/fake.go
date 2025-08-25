package fake

import (
	"context"
	"math/rand"
	"time"

	"github.com/jstaud/edgelink-go/pkg/models"
)

// Client simulates a device for testing
// It "implements" the Driver interface by having the required methods
type Client struct {
	spec models.DeviceSpec  // Configuration for this fake device
	rnd  *rand.Rand        // Random number generator for realistic data
}

// New creates a new fake device client
// This is a "constructor function" - common Go pattern
func New(spec models.DeviceSpec) *Client {
	return &Client{
		spec: spec,
		// Create a random number generator with current time as seed
		rnd:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Poll simulates reading sensor data
// This implements the Driver interface
func (c *Client) Poll(ctx context.Context) (models.Reading, error) {
	// Create empty map to hold our fake sensor readings
	metrics := map[string]float64{}
	
	// Generate fake data for each metric this device reports
	for _, name := range c.spec.MetricNames {
		base := 0.0
		
		// Generate realistic base values for different sensor types
		switch name {
		case "temp_c":
			base = 60 + c.rnd.Float64()*20  // 60-80°C
		case "pressure_bar":
			base = 6 + c.rnd.Float64()*0.5  // 6-6.5 bar
		case "vibration_rms":
			base = 0.2 + c.rnd.Float64()*0.05  // 0.2-0.25 RMS
		default:
			base = c.rnd.Float64()*100  // 0-100 for unknown metrics
		}
		
		// Add jitter (noise) to make it realistic
		// jitter is a percentage (0.05 = 5%)
		if c.spec.Jitter > 0 {
			// Random value between -1 and +1, multiplied by jitter percentage
			jitterMultiplier := 1 + (c.spec.Jitter * (c.rnd.Float64()*2 - 1))
			metrics[name] = base * jitterMultiplier
		} else {
			metrics[name] = base
		}
	}
	
	// Create the reading with current timestamp
	r := models.Reading{
		DeviceID:  c.spec.ID,
		Timestamp: time.Now().UTC(),  // Always use UTC for consistency
		Metrics:   metrics,
		Quality:   "good",  // Fake devices always have good quality
	}
	
	return r, nil  // Return reading and no error
}

// Close cleans up resources (fake device has nothing to clean up)
func (c *Client) Close() error {
	return nil
}
