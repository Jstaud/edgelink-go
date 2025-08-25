package device

import (
	"context"
	"github.com/jstaud/edgelink-go/pkg/models"
)

// Driver defines what all device drivers must implement
// This is an INTERFACE - like a contract that says "you must have these methods"
type Driver interface {
	// Poll reads data from the device
	// Takes a context (for timeouts/cancellation) and returns Reading + error
	Poll(ctx context.Context) (models.Reading, error)
	
	// Close cleans up resources when we're done
	Close() error
}
