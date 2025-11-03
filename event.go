package hamqtt

const (
	EventButton   = "button"
	EventDoorbell = "doorbell"
	EventMotion   = "motion"
)

// EventConfig provides configuration for event entities.
type EventConfig struct {

	// ID provides the unique ID of the binary sensor.
	ID string

	// Name provides the name of the binary sensor.
	Name string

	// DeviceClass categorizes the type of event.
	DeviceClass string

	// EventTypes is a list of valid events (i.e. "press", "hold", etc.).
	EventTypes []string
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
func (c *Conn) Event(cfg *EventConfig) (*Event, error) {
	stateTopic := c.stateTopic(cfg.ID)
	if err := c.publishCfg(
		c.cfgTopic(cfg.ID, "event"),
		map[string]any{
			"device":       c.device,
			"platform":     "event",
			"unique_id":    cfg.ID,
			"name":         cfg.Name,
			"device_class": cfg.DeviceClass,
			"event_types":  cfg.EventTypes,
			"state_topic":  stateTopic,
		},
	); err != nil {
		return nil, err
	}
	return &Event{
		conn:       c,
		stateTopic: stateTopic,
	}, nil
}
