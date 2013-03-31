package huego

// Change describes a set of pending changes to the state of a light.
type Change struct {
	hub    *Hub
	path   string
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
	_, err := c.hub.Put(c.path, c.params, nil)
	return err
}
