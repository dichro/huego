package huego

import (
	"fmt"
)

// Light describes various parameters of a light, as returned from the hub.
type Light struct {
	hub             *Hub
	id              string
	State           State  `json:"state"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Model           string `json:"modelid"`
	SoftwareVersion string `json:"swversion"`
	// what's a pointsymbol?
}

// State describes the current state of the light, as returned from the hub.
type State struct {
	On         bool      `json:"on"`
	Brightness int       `json:"bri"`
	Hue        int       `json:"hue"`
	Saturation int       `json:"sat"`
	CIE        []float32 `json:"xy"`
	Mireds     int       `json:"ct"`
	Alert      string    `json:"alert"`
	Effect     string    `json:"effect"`
	ColourMode string    `json:"colormode"`
	Reachable  bool      `json:"reachable"`
}

func (l Light) String() string {
	state := "off"
	if l.State.On {
		state = "on"
	}
	return fmt.Sprintf("%s: %s", l.Name, state)
}

func (l *Light) Switch() *Change {
	return &Change{
		hub:    l.hub,
		path:   fmt.Sprintf("lights/%s/state", l.id),
		params: make(map[string]interface{}),
	}
}
