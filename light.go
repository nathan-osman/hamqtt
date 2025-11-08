package hamqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// LightConfig provides configuration for light entities.
type LightConfig struct {

	// State indicates the initial state of the light.
	State bool `json:"-"`

	// OnCallback is invoked when the light is turned on. Returning true will
	// cause a change to the light's state.
	OnCallback func() bool `json:"-"`

	// OffCallback is invoked when the light is turned off. Returning true
	// will cause a change to the light's state.
	OffCallback func() bool `json:"-"`
}

type hamqttLight struct {
	*EntityConfig
	*LightConfig
	Device       *hamqttDevice `json:"device"`
	Platform     string        `json:"platform"`
	CommandTopic string        `json:"command_topic"`
	StateTopic   string        `json:"state_topic"`
}

// Light creates a new light entity with the provided configuration.
func (c *Conn) Light(
	entityCfg *EntityConfig,
	cfg *LightConfig,
) error {
	var (
		cmdTopic   = c.cmdTopic(entityCfg.ID)
		stateTopic = c.stateTopic(entityCfg.ID)
	)
	if err := c.publishStateBool(stateTopic, cfg.State); err != nil {
		return err
	}
	if err := c.publishCfg(
		c.cfgTopic(entityCfg.ID, "light"),
		&hamqttLight{
			EntityConfig: entityCfg,
			LightConfig:  cfg,
			Device:       c.device,
			Platform:     "light",
			CommandTopic: cmdTopic,
			StateTopic:   stateTopic,
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
					c.publishStateBool(stateTopic, true)
				}
			case "OFF":
				if cfg.OffCallback() {
					c.publishStateBool(stateTopic, false)
				}
			}
		},
	); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
