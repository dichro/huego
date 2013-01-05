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
}

type Light struct {
	State State `json:"state"`
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