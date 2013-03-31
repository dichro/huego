package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Action Action   `json:"action"`
	Lights []string `json:"lights"`
	Name   string   `json:"name"`
}

type Status struct {
	Lights map[string]Light `json:"lights"`
	Groups map[string]Group `json:"groups"`
	Config Config           `json:"config"`
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

func (h *Hub) ChangeLight(light string) *Change {
	return &Change{
		hub:    h,
		id:     light,
		params: make(map[string]interface{}),
	}
}
