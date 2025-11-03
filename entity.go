package hamqtt

type EntityConfig struct {

	// ID provides the unique ID of the entity.
	ID string `json:"unique_id"`

	// Name provides the name of the entity.
	Name string `json:"name"`
}
