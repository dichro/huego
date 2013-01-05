package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type State struct {
	On bool `json:"on"`
	Brightness int `json:"bri"`
	Hue int `json:"hue"`
	Saturation int `json:"sat"`
	CIE []float32 `json:"xy"`
	Mireds int `json:"ct"`
	Alert string `json:"alert"`
	Effect string `json:"effect"`
	ColourMode string `json:"colormode"`
	Reachable bool `json:"reachable"`
}

type Light struct {
	State State `json:"state"`
	Type string `json:"type"`
	Name string `json:"name"`
	Model string `json:"modelid"`
	SoftwareVersion string `json:"swversion"`
	// what's a pointsymbol?
}

func (l Light) String() string {
	state := "off"
	if l.State.On {
		state = "on"
	}
	return fmt.Sprintf("%s: %s", l.Name, state)
}

type ACL struct {
	Last string `json:"last use date"`
	Created string `json:"create date"`
	Name string `json:"name"`
}

type Update struct {
	State int `json:"updatestate"`
	URL string `json:"url"`
	Text string `json:"text"`
	Notify bool `json:"notify"`
}

type Config struct {
	Name string `json:"name"`
	MAC string `json:"mac"`
	DHCP bool `json:"dhcp"`
	IP string `json:"ipaddress"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	ProxyIP string `json:"proxyaddress"`
	ProxyPort int `json:"proxyport"`
	UTC string
	ACL map[string]ACL `json:"whitelist"`
	SoftwareVersion string `json:"swversion"`
	Update Update `json:"swupdate"`
	LinkButton bool `json:"linkbutton"`
	PortalServices bool `json:"portalservices"`
}

type Status struct {
	Lights map[string]Light `json:"lights"`
	// TODO(dichro): groups
	Config Config `json:"config"`
	// TODO(dichro): schedules
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
	return status, err
}


func (h *Hub) SetLightState(light string, state State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/api/%s/lights/%s/state", h.Address, h.Username, light), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err, resp)
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	ret := make(map[string]interface{})
	err = dec.Decode(ret)
	fmt.Println(ret)
	return err
}
