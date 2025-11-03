package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/nathan-osman/hamqtt"
)

func main() {
	var (
		addr     = flag.String("addr", "", "MQTT broker address")
		username = flag.String("username", "", "Username for authentication")
		password = flag.String("password", "", "Password for authentication")
	)

	// Parse CLI arguments
	flag.Parse()

	// Connect to the MQTT broker
	c, err := hamqtt.New(&hamqtt.Config{
		Addr:     *addr,
		Username: *username,
		Password: *password,
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Create a thermometer
	s, err := c.Sensor(&hamqtt.SensorConfig{
		ID:                        "mythermometer",
		Name:                      "My Thermometer",
		DeviceClass:               hamqtt.SensorTemperature,
		UnitOfMeasurement:         hamqtt.SensorDegreesCelsius,
		SuggestedDisplayPrecision: 1,
	})
	if err != nil {
		panic(err)
	}

	// Read interactive commands from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("New value for temperature in Celsius? (EOF to quit)")
		if !scanner.Scan() {
			return
		}
		s.Set(scanner.Text())
	}
}
