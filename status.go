package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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