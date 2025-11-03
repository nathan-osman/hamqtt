package hamqtt

const (
	BinarySensorBattery         = "battery"
	BinarySensorBatteryCharging = "battery_charging"
	BinarySensorCarbonMonoxide  = "carbon_monoxide"
	BinarySensorCold            = "cold"
	BinarySensorConnectivity    = "connectivity"
	BinarySensorDoor            = "door"
	BinarySensorGarageDpor      = "garage_door"
	BinarySensorGas             = "gas"
	BinarySensorHeat            = "heat"
	BinarySensorLight           = "light"
	BinarySensorLock            = "lock"
	BinarySensorMoisture        = "moisture"
	BinarySensorMotion          = "motion"
	BinarySensorMoving          = "moving"
	BinarySensorOccupancy       = "occupancy"
	BinarySensorOpening         = "opening"
	BinarySensorPlug            = "plug"
	BinarySensorPower           = "power"
	BinarySensorPresence        = "presence"
	BinarySensorProblem         = "problem"
	BinarySensorRunning         = "running"
	BinarySensorSafety          = "safety"
	BinarySensorSmoke           = "smoke"
	BinarySensorSound           = "sound"
	BinarySensorTamper          = "tamper"
	BinarySensorUpdate          = "update"
	BinarySensorVibration       = "vibration"
	BinarySensorWindow          = "window"
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

// Set updates the binary sensor's value.
func (b *BinarySensor) Set(state bool) error {
	return b.conn.publishStateBool(b.stateTopic, state)
}

// BinarySensor creates a new entity that represents a binary sensor.
func (c *Conn) BinarySensor(cfg *BinarySensorConfig) (*BinarySensor, error) {
	stateTopic := c.stateTopic(cfg.ID)
	if err := c.publishStateBool(stateTopic, cfg.State); err != nil {
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
