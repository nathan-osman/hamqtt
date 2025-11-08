package hamqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NotifyConfig provides configuration for notify entities.
type NotifyConfig struct {

	// NotifyCallback is invoked when a new notification is received.
	NotifyCallback func(payload string) `json:"-"`
}

type hamqttNotify struct {
	*EntityConfig
	*NotifyConfig
	Device       *hamqttDevice `json:"device"`
	CommandTopic string        `json:"command_topic"`
}

// Notify creates a new notify entity with the provided configuration.
func (c *Conn) Notify(
	entityCfg *EntityConfig,
	cfg *NotifyConfig,
) error {
	cmdTopic := c.cmdTopic(entityCfg.ID, "notify")
	if err := c.publishCfg(
		c.cfgTopic(entityCfg.ID, "notify"),
		&hamqttNotify{
			EntityConfig: entityCfg,
			NotifyConfig: cfg,
			Device:       c.device,
			CommandTopic: cmdTopic,
		},
	); err != nil {
		return err
	}
	if t := c.client.Subscribe(
		cmdTopic,
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			cfg.NotifyCallback(string(msg.Payload()))
		},
	); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
