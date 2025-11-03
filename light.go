package hamqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

	// OffCallback is invoked when the light is turned off. Returning true
	// will cause a change to the light's state.
	OffCallback BoolCallback
}

// Light creates a new light entity with the provided configuration.
func (c *Conn) Light(cfg *LightConfig) error {
	var (
		cmdTopic   = c.cmdTopic(cfg.ID)
		stateTopic = c.stateTopic(cfg.ID)
	)
	if err := c.publishState(stateTopic, cfg.State); err != nil {
		return err
	}
	if err := c.publishCfg(
		c.cfgTopic(cfg.ID, "light"),
		map[string]any{
			"device":        c.device,
			"platform":      "light",
			"unique_id":     cfg.ID,
			"name":          cfg.Name,
			"command_topic": cmdTopic,
			"state_topic":   stateTopic,
		},
	); err != nil {
		return err
	}
	if t := c.client.Subscribe(
		cmdTopic,
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			switch string(msg.Payload()) {
			case "ON":
				if cfg.OnCallback() {
					c.publishState(stateTopic, true)
				}
			case "OFF":
				if cfg.OffCallback() {
					c.publishState(stateTopic, false)
				}
			}
		},
	); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
