package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ACL struct {
	Last    string `json:"last use date"`
	Created string `json:"create date"`
	Name    string `json:"name"`
}

type Update struct {
	State  int    `json:"updatestate"`
	URL    string `json:"url"`
	Text   string `json:"text"`
	Notify bool   `json:"notify"`
}

type Config struct {
	Name            string `json:"name"`
	MAC             string `json:"mac"`
	DHCP            bool   `json:"dhcp"`
	IP              string `json:"ipaddress"`
	Netmask         string `json:"netmask"`
	Gateway         string `json:"gateway"`
	ProxyIP         string `json:"proxyaddress"`
	ProxyPort       int    `json:"proxyport"`
	UTC             string
	ACL             map[string]ACL `json:"whitelist"`
	SoftwareVersion string         `json:"swversion"`
	Update          Update         `json:"swupdate"`
	LinkButton      bool           `json:"linkbutton"`
	PortalServices  bool           `json:"portalservices"`
}

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

type Status struct {
	hub *Hub
	Lights map[string]*Light `json:"lights"`
	Groups map[string]*Group `json:"groups"`
	Config Config           `json:"config"`
	// TODO(dichro): schedules
}

// assign updates all useful state in Status and its children and must
// be called after json parsing.
func (s *Status) assign(h *Hub) {
	s.hub = h
	for id, group := range s.Groups {
		group.status = s
		group.id = id
	}
}

func (h *Hub) Status() (*Status, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/%s/", h.Address, h.Username))
	if err != nil {
		log.Println(err, resp)
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	status := &Status{}
	err = dec.Decode(status)
	status.assign(h)
	return status, err
}

func (h *Hub) ChangeLight(light string) *Change {
	return &Change{
		hub:    h,
		id:     light,
		params: make(map[string]interface{}),
	}
}
