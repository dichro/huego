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
	status *Status
	id string
	Action Action   `json:"action"`
	Lights []string `json:"lights"`
	Name   string   `json:"name"`
}

func (g *Group) Switch() *GroupChange {
	return &GroupChange{
		group: g,
		params: make(map[string]interface{}),
	}
}

type GroupChange struct {
	group *Group
	params map[string]interface{}
}

// State requests that the light be set to the requested state.
func (c *GroupChange) State(on bool) *GroupChange {
	c.params["on"] = on
	return c
}

// Send dispatches all the requested changes to the light.
func (c *GroupChange) Send() error {
	_, err := c.group.status.hub.Put(fmt.Sprintf("groups/%s/action", c.group.id), c.params, nil)
	return err
}
