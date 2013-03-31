package huego

import (
	"fmt"
)

type Action struct {
	On         bool      `json:"on"`
	Brightness int       `json:"bri"`
	Hue        int       `json:"hue"`
	Saturation int       `json:"sat"`
	CIE        []float32 `json:"xy"`
	Mireds     int       `json:"ct"`
	Effect     string    `json:"effect"`
	ColourMode string    `json:"colormode"`
}

type Group struct {
	hub    *Hub
	id     string
	Action Action   `json:"action"`
	Lights []string `json:"lights"`
	Name   string   `json:"name"`
}

func (g *Group) Switch() *Change {
	return &Change{
		hub:    g.hub,
		path:   fmt.Sprintf("groups/%s/action", g.id),
		params: make(map[string]interface{}),
	}
}
