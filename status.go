package huego

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	SoftwareVer string `json:"swversion"`
	// what's a pointsymbol?
}

type Status struct {
	Lights map[string]Light `json:"lights"`
}
	
func (h *Hub) Status() {
	if resp, err := http.Get(fmt.Sprintf("http://%s/api/%s/", h.Address, h.Username)); err != nil {
		log.Println(err, resp)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		sbody := string(body)
		fmt.Println(sbody)
		dec := json.NewDecoder(strings.NewReader(sbody))
		status := Status{}
		err = dec.Decode(&status)
		fmt.Println(err, status)
	}
}