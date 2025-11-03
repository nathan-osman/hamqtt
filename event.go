package hamqtt

const (
	EventButton   = "button"
	EventDoorbell = "doorbell"
	EventMotion   = "motion"
)

// EventConfig provides configuration for event entities.
type EventConfig struct {

	// DeviceClass categorizes the type of event.
	DeviceClass string `json:"device_class,omitempty"`

	// EventTypes is a list of valid events (i.e. "press", "hold", etc.).
	EventTypes []string `json:"event_types,omitempty"`
}

type hamqttEvent struct {
	*EntityConfig
	*EventConfig
	Device     *hamqttDevice `json:"device"`
	Platform   string        `json:"platform"`
	StateTopic string        `json:"state_topic"`
}

// Event provides methods for sending events.
type Event struct {
	conn       *Conn
	stateTopic string
}

// Send sends the specified event.
func (e *Event) Send(eventType string) error {
	return e.conn.publishStateJSON(
		e.stateTopic,
		map[string]any{
			"event_type": eventType,
		},
	)
}

// Event creates a new event entity with the provided configuration.
func (c *Conn) Event(
	entityCfg *EntityConfig,
	cfg *EventConfig,
) (*Event, error) {
	stateTopic := c.stateTopic(entityCfg.ID)
	if err := c.publishCfg(
		c.cfgTopic(entityCfg.ID, "event"),
		&hamqttEvent{
			EntityConfig: entityCfg,
			EventConfig:  cfg,
			Device:       c.device,
			Platform:     "event",
			StateTopic:   stateTopic,
		},
	); err != nil {
		return nil, err
	}
	return &Event{
		conn:       c,
		stateTopic: stateTopic,
	}, nil
}
