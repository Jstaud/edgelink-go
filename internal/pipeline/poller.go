package pipeline

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/jstaud/edgelink-go/internal/device"
	"github.com/jstaud/edgelink-go/pkg/models"
)

// Publisher interface - anything that can publish readings
// This lets us swap MQTT for Kafka, Redis, etc.
type Publisher interface {
	Publish(ctx context.Context, r models.Reading) error
}

// Cache interface - anything that can store latest readings
type Cache interface {
	Set(r models.Reading)
	Latest(id string) (models.Reading, bool)
}

// RunPoller polls a device on a schedule and publishes readings
// This runs in its own goroutine (like a background thread)
func RunPoller(ctx context.Context, log *zap.Logger, drv device.Driver, spec models.DeviceSpec, pub Publisher, cache Cache) error {
	// Create a ticker that fires every poll_every duration
	ticker := time.NewTicker(spec.PollEvery)
	defer ticker.Stop()  // Clean up ticker when function exits

	// Infinite loop until context is cancelled
	for {
		select {
		case <-ctx.Done():
			// Context was cancelled (shutdown signal), exit gracefully
			return ctx.Err()
			
		case <-ticker.C:
			// Time to poll! The ticker fired
			
			// Give each poll operation half the poll interval as timeout budget
			// This prevents one slow poll from blocking the next one
			cctx, cancel := context.WithTimeout(ctx, spec.PollEvery/2)
			
			// Poll the device
			r, err := drv.Poll(cctx)
			cancel()  // Always cancel context to free resources
			
			if err != nil {
				// Only log unexpected errors, not cancellation/timeout
				if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
					log.Warn("poll error", 
						zap.String("device", spec.ID), 
						zap.Error(err))
				}
				continue  // Skip to next poll attempt
			}
			
			// Store reading in cache for REST API
			cache.Set(r)
			
			// Publish reading to MQTT
			if err := pub.Publish(ctx, r); err != nil {
				log.Warn("publish error", 
					zap.String("device", spec.ID), 
					zap.Error(err))
				// Note: we don't return here - keep polling even if publish fails
			}
		}
	}
}
