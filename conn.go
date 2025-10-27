package hamqtt

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type hamqttDevice struct {
	IDs          string `json:"identifiers"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
}

// Conn maintains a connection to an MQTT broker.
type Conn struct {
	client          mqtt.Client
	discoveryPrefix string
	id              string
	device          *hamqttDevice
}

// New creates a new Conn instance with the provided configuration.
func New(cfg *Config) (*Conn, error) {

	// If DiscoveryPrefix is empty, use the default
	discoveryPrefix := cfg.DiscoveryPrefix
	if discoveryPrefix == "" {
		discoveryPrefix = "homeassistant"
	}

	// If ID and / or Name was not provided, use the hostname
	var (
		id   = cfg.ID
		name = cfg.Name
	)
	if id == "" {
		v, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		id = v
	}
	if name == "" {
		name = id
	}

	// Connect to the MQTT broker
	c := mqtt.NewClient(
		mqtt.NewClientOptions().
			AddBroker(fmt.Sprintf("tcp://%s", cfg.Addr)).
			SetUsername(cfg.Username).
			SetPassword(cfg.Password).
			SetClientID(id).
			SetResumeSubs(true),
	)

	// Ensure the connection was successful
	if t := c.Connect(); t.Wait() && t.Error() != nil {
		return nil, t.Error()
	}

	// Return the Conn
	return &Conn{
		client:          c,
		discoveryPrefix: discoveryPrefix,
		id:              id,
		device: &hamqttDevice{
			IDs:          id,
			Name:         name,
			Manufacturer: cfg.Manufacturer,
			Model:        cfg.Model,
		},
	}, nil
}

// Close shuts down the connection.
func (c *Conn) Close() {
	c.client.Disconnect(1000)
}
