package hamqtt

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// BoolCallback is a callback that returns a bool indicating success.
type BoolCallback func() bool

// LightConfig provides configuration for light entities.
type LightConfig struct {

	// ID provides the unique ID of the light.
	ID string

	// Name provides the name of the light.
	Name string

	// State indicates the initial state of the light.
	State bool

	// OnCallback is invoked when the light is turned on. Returning true will
	// cause a change to the light's state.
	OnCallback BoolCallback

	// OffCallback is invoked when
	OffCallback BoolCallback
}

func (c *Conn) lightPublishState(topic string, state bool) error {
	payload := "OFF"
	if state {
		payload = "ON"
	}
	if t := c.client.Publish(topic, 0, true, payload); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

// Light creates a new light entity with the provided configuration.
func (c *Conn) Light(cfg *LightConfig) error {
	var (
		cfgTopic = fmt.Sprintf(
			"%s/light/%s/config",
			c.discoveryPrefix,
			cfg.ID,
		)
		cmdTopic = fmt.Sprintf(
			"%s/%s/switch",
			c.id,
			cfg.ID,
		)
		stateTopic = fmt.Sprintf(
			"%s/%s/state",
			c.id,
			cfg.ID,
		)
		payload = map[string]any{
			"device":        c.device,
			"platform":      "light",
			"unique_id":     cfg.ID,
			"name":          cfg.Name,
			"command_topic": cmdTopic,
			"state_topic":   stateTopic,
		}
	)
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := c.lightPublishState(stateTopic, cfg.State); err != nil {
		return err
	}
	if t := c.client.Publish(cfgTopic, 0, true, b); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	if t := c.client.Subscribe(
		cmdTopic,
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			switch string(msg.Payload()) {
			case "ON":
				if cfg.OnCallback() {
					c.lightPublishState(stateTopic, true)
				}
			case "OFF":
				if cfg.OffCallback() {
					c.lightPublishState(stateTopic, false)
				}
			}
		},
	); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
