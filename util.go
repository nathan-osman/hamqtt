package hamqtt

import (
	"encoding/json"
	"fmt"
)

func (c *Conn) cfgTopic(id, entityType string) string {
	return fmt.Sprintf(
		"%s/%s/%s/config",
		c.discoveryPrefix,
		entityType,
		id,
	)
}

func (c *Conn) cmdTopic(id string) string {
	return fmt.Sprintf(
		"%s/%s/switch",
		c.id,
		id,
	)
}

func (c *Conn) stateTopic(id string) string {
	return fmt.Sprintf(
		"%s/%s/state",
		c.id,
		id,
	)
}

func (c *Conn) publishCfg(topic string, payload map[string]any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if t := c.client.Publish(topic, 0, true, b); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

func (c *Conn) publishState(topic string, state bool) error {
	payload := "OFF"
	if state {
		payload = "ON"
	}
	if t := c.client.Publish(topic, 0, true, payload); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
