package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Light describes various parameters of a light, as returned from the hub.
type Light struct {
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

// TODO(dichro): what's the easiest way to propagate the hub and id
// all the way down to here? Magic with UnmarshalJSON?
func (l Light) change() *Change {
	return &Change{
		params: make(map[string]interface{}),
	}
}

// Change describes a set of pending changes to the state of a light.
type Change struct {
	hub    *Hub
	id     string
	params map[string]interface{}
}

// On updates this Change to request that the light be turned on.
func (c *Change) On() *Change {
	c.params["on"] = true
	return c
}

// Off updates this Change to request that the light be turned off.
func (c *Change) Off() *Change {
	c.params["on"] = false
	return c
}

// State updates this Change to request that the light be set to the
// requested state.
func (c *Change) State(on bool) *Change {
	c.params["on"] = on
	return c
}

// Transition sets the time that this change should be applied over.
func (c *Change) Transition(centiSeconds int) *Change {
	c.params["transitiontime"] = centiSeconds
	return c
}

// Temperature sets the requested colour temperature.
func (c *Change) Temperature(temp int) *Change {
	c.params["ct"] = temp
	return c
}

// Brightness sets the requested brightness.
func (c *Change) Brightness(bri int) *Change {
	c.params["bri"] = bri
	return c
}

// Colour sets the requested colour.
func (c *Change) Colour(hue, saturation int) *Change {
	c.params["hue"], c.params["sat"] = hue, saturation
	return c
}

// Send dispatches all the requested changes to the light.
func (c *Change) Send() error {
	data, err := json.Marshal(c.params)
	if err != nil {
		return err
	}
	log.Println("request", string(data))
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/api/%s/lights/%s/state", c.hub.Address, c.hub.Username, c.id), strings.NewReader(string(data)))
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
