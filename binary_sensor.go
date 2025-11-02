package hamqtt

const (
	Battery         = "battery"
	BatteryCharging = "battery_charging"
	CarbonMonoxide  = "carbon_monoxide"
	Cold            = "cold"
	Connectivity    = "connectivity"
	Door            = "door"
	GarageDpor      = "garage_door"
	Gas             = "gas"
	Heat            = "heat"
	Light           = "light"
	Lock            = "lock"
	Moisture        = "moisture"
	Motion          = "motion"
	Moving          = "moving"
	Occupancy       = "occupancy"
	Opening         = "opening"
	Plug            = "plug"
	Power           = "power"
	Presence        = "presence"
	Problem         = "problem"
	Running         = "running"
	Safety          = "safety"
	Smoke           = "smoke"
	Sound           = "sound"
	Tamper          = "tamper"
	Update          = "update"
	Vibration       = "vibration"
	Window          = "window"
)

// BinarySensorConfig provides configuration for binary sensors.
type BinarySensorConfig struct {

	// ID provides the unique ID of the binary sensor.
	ID string

	// Name provides the name of the binary sensor.
	Name string

	// State indicates the initial state of the binary sensor.
	State bool

	// DeviceClass categorizes the type of data reported by the sensor.
	DeviceClass string
}

// BinarySensor provides methods for indicating changes to the sensor.
type BinarySensor struct {
	conn       *Conn
	stateTopic string
}

func (b *BinarySensor) Set(state bool) error {
	return b.conn.publishState(b.stateTopic, state)
}

// BinarySensor creates a new entity that represents a binary sensor.
func (c *Conn) BinarySensor(cfg *BinarySensorConfig) (*BinarySensor, error) {
	stateTopic := c.stateTopic(cfg.ID)
	if err := c.publishState(stateTopic, cfg.State); err != nil {
		return nil, err
	}
	if err := c.publishCfg(
		c.cfgTopic(cfg.ID, "binary_sensor"),
		map[string]any{
			"device":       c.device,
			"platform":     "binary_sensor",
			"unique_id":    cfg.ID,
			"name":         cfg.Name,
			"device_class": cfg.DeviceClass,
			"state_topic":  stateTopic,
		},
	); err != nil {
		return nil, err
	}
	return &BinarySensor{
		conn:       c,
		stateTopic: stateTopic,
	}, nil
}
