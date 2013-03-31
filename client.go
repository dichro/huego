/*
 Package huego provides a client for Philips Hue smart lightbulbs.
*/
package huego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Hub struct {
	Username, Address string
}

func (h *Hub) Put(path string, request, response interface{}) (interface{}, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s/api/%s/%s", h.Address, h.Username, path)
	log.Printf("hub PUT [%s]: %q %q", h, uri, request)
	req, err := http.NewRequest("PUT", uri, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err, resp)
		return nil, err
	}
	defer resp.Body.Close()
	if response == nil {
		response = make(map[string]interface{})
	}
	// TODO(dichro): what does this actually return?
	dec := json.NewDecoder(resp.Body)
	defer log.Printf("response %q", response)
	return response, dec.Decode(response)
}
