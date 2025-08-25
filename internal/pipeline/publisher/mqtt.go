package publisher

import (
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jstaud/edgelink-go/pkg/models"
)

// MQTTPublisher sends readings to an MQTT broker
type MQTTPublisher struct {
	c         mqtt.Client  // The MQTT client connection
	baseTopic string       // Base topic prefix (e.g., "plant/readings")
}

// NewMQTT creates a new MQTT publisher
func NewMQTT(brokerURL, clientID, baseTopic string) *MQTTPublisher {
	// Configure MQTT client options
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).                    // Where to connect
		SetClientID(clientID).                   // Unique ID for this client
		SetAutoReconnect(true).                  // Reconnect if connection drops
		SetConnectRetry(true).                   // Retry if initial connection fails
		SetConnectTimeout(3 * time.Second)      // Give up connecting after 3 seconds
	
	// Create and connect the client
	c := mqtt.NewClient(opts)
	token := c.Connect()  // This returns a "token" that tracks the operation
	token.Wait()          // Wait for connection to complete
	
	return &MQTTPublisher{
		c:         c,
		baseTopic: baseTopic,
	}
}

// Publish sends a reading to MQTT
func (p *MQTTPublisher) Publish(ctx context.Context, r models.Reading) error {
	// Convert reading to JSON
	payload, _ := json.Marshal(r)
	
	// Create topic: baseTopic + "/" + deviceID
	// e.g., "plant/readings/press-1"
	topic := p.baseTopic + "/" + r.DeviceID
	
	// Publish with QoS 1 (at least once delivery), not retained
	token := p.c.Publish(topic, 1, false, payload)
	token.Wait()  // Wait for publish to complete
	
	return token.Error()  // Return any error that occurred
}

// Close disconnects from MQTT broker
func (p *MQTTPublisher) Close() {
	if p.c.IsConnected() {
		p.c.Disconnect(250)  // Wait up to 250ms for graceful disconnect
	}
}
