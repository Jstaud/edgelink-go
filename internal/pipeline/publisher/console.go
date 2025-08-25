package publisher

import (
	"context"
	"fmt"
	"github.com/jstaud/edgelink-go/pkg/models"
)

// ConsolePublisher just prints readings to console instead of MQTT
type ConsolePublisher struct {
	baseTopic string
}

// NewConsole creates a console publisher for testing
func NewConsole(baseTopic string) *ConsolePublisher {
	return &ConsolePublisher{baseTopic: baseTopic}
}

// Publish prints the reading instead of sending to MQTT
func (p *ConsolePublisher) Publish(ctx context.Context, r models.Reading) error {
	fmt.Printf("[%s] %s/%s: %+v\n", r.Timestamp.Format("15:04:05"), p.baseTopic, r.DeviceID, r.Metrics)
	return nil
}

// Close is a no-op for console publisher
func (p *ConsolePublisher) Close() {
	// Nothing to close
}
