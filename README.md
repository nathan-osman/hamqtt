## hamqtt

[![Go Reference](https://pkg.go.dev/badge/github.com/nathan-osman/hamqtt.svg)](https://pkg.go.dev/github.com/nathan-osman/hamqtt)
[![MIT License](https://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](https://opensource.org/licenses/MIT)

This package aims to provide an easy interface for exposing entities to [Home Assistant](https://www.home-assistant.io/) via [MQTT](https://mqtt.org/).

### Installation

Adding hamqtt to your application is as easy as:

```golang
import "github.com/nathan-osman/hamqtt"
```

### Usage

To expose entities to Home Assistant, you must first create a connection (`Conn`):

```golang
c, err := hamqtt.New(&hamqtt.Config{
    Addr:     "1.2.3.4:1883",
    Username: "myusername",
    Password: "password123",
})
if err != nil {
    panic(err)
}
defer c.Close()
```

From there, you can create entities by calling their respective member function. For example, to create a light:

```golang
c.Light(
    &hamqtt.EntityConfig{
        ID:   "mylight",
        Name: "My Light",
    },
    &hamqtt.LightConfig{
        OnCallback: func() bool {
            fmt.Println("Light turned on")
            return true
        },
        OffCallback: func() bool {
            fmt.Println("Light turned off")
            return true
        },
    },
)
```

This will create a corresponding entity for your light that you can then control directly from the Home Assistant UI:

![Screenshot of light control in Home Assistant](https://github.com/nathan-osman/hamqtt/blob/main/dist/example-light.png?raw=true)

Toggling the switch should cause the appropriate message to be written to your application's console.
