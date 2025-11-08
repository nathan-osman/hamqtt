package hamqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	ButtonIdentify = "identify"
	ButtonRestart  = "restart"
	ButtonUpdate   = "update"
)

// ButtonConfig provides configuration for button entities.
type ButtonConfig struct {

	// DeviceClass categorizes the button.
	DeviceClass string `json:"device_class,omitempty"`

	// PressCallback is invoked when the button is pressed.
	PressCallback func() `json:"-"`
}

type hamqttButton struct {
	*EntityConfig
	*ButtonConfig
	Device       *hamqttDevice `json:"device"`
	Platform     string        `json:"platform"`
	CommandTopic string        `json:"command_topic"`
}

// Button creates a new button entity with the provided configuration.
func (c *Conn) Button(
	entityCfg *EntityConfig,
	cfg *ButtonConfig,
) error {
	cmdTopic := c.cmdTopic(entityCfg.ID)
	if err := c.publishCfg(
		c.cfgTopic(entityCfg.ID, "button"),
		&hamqttButton{
			EntityConfig: entityCfg,
			ButtonConfig: cfg,
			Device:       c.device,
			Platform:     "button",
			CommandTopic: cmdTopic,
		},
	); err != nil {
		return err
	}
	if t := c.client.Subscribe(
		cmdTopic,
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			if string(msg.Payload()) == "PRESS" {
				cfg.PressCallback()
			}
		},
	); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
