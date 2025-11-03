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

	// ID provides the unique ID of the button.
	ID string

	// Name provides the name of the button.
	Name string

	// DeviceClass categorizes the button.
	DeviceClass string

	// PressCallback is invoked when the button is pressed.
	PressCallback EmptyCallback
}

// Button creates a new button entity with the provided configuration.
func (c *Conn) Button(cfg *ButtonConfig) error {
	cmdTopic := c.cmdTopic(cfg.ID)
	if err := c.publishCfg(
		c.cfgTopic(cfg.ID, "button"),
		map[string]any{
			"device":        c.device,
			"platform":      "button",
			"unique_id":     cfg.ID,
			"name":          cfg.Name,
			"device_class":  cfg.DeviceClass,
			"command_topic": cmdTopic,
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
