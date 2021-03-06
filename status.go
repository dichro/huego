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

type Status struct {
	hub    *Hub
	Lights map[string]*Light `json:"lights"`
	Groups map[string]*Group `json:"groups"`
	Config Config            `json:"config"`
	// TODO(dichro): schedules
}

type Switchable interface {
	Switch() *Change
}

func (s *Status) FindSwitchable(name string) Switchable {
	for _, group := range s.Groups {
		if group.Name == name {
			return group
		}
	}
	for _, light := range s.Lights {
		if light.Name == name {
			return light
		}
	}
	return nil
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
	status.hub = h
	for id, group := range status.Groups {
		group.hub = h
		group.id = id
	}
	for id, light := range status.Lights {
		light.hub = h
		light.id = id
	}
	return status, err
}
