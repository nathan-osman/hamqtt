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

func (c *Conn) cmdTopic(id, entityType string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		c.id,
		id,
		entityType,
	)
}

func (c *Conn) stateTopic(id string) string {
	return fmt.Sprintf(
		"%s/%s/state",
		c.id,
		id,
	)
}

func (c *Conn) publishCfg(topic string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if t := c.client.Publish(topic, 0, true, b); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

func (c *Conn) publishState(topic string, payload string) error {
	if t := c.client.Publish(topic, 0, true, payload); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

func (c *Conn) publishStateBool(topic string, state bool) error {
	payload := "OFF"
	if state {
		payload = "ON"
	}
	return c.publishState(topic, payload)
}

func (c *Conn) publishStateJSON(topic string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.publishState(topic, string(b))
}
