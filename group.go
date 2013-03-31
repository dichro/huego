package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	data, err := json.Marshal(c.params)
	if err != nil {
		return err
	}
	hub := c.group.status.hub
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/api/%s/groups/%s/action", hub.Address, hub.Username, c.group.id), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err, resp)
		return err
	}
	defer resp.Body.Close()
	// TODO(dichro): what does this actually return?
	dec := json.NewDecoder(resp.Body)
	ret := make(map[string]interface{})
	log.Println("response", ret)
	return dec.Decode(ret)
}
