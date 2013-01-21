package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dichro/huego"
)

var (
	ok = []byte("OK")
	username = flag.String("username", "", "username for Hue hub")
	address  = flag.String("address", "", "address of Hue hub")
	re = regexp.MustCompile("^/[a-z]+/([A-Za-z0-9' ]+)/([.0-9]+)$")
)


func setBrightness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	groups := re.FindStringSubmatch(req.URL.Path)
	if len(groups) != 3 {
		w.Write([]byte("bad"))
		return
	}
	name := groups[1]
	bri, err := strconv.ParseFloat(groups[2], 64)
	if err != nil {
		w.Write([]byte("bad"))
		return
	}
	hub := &huego.Hub{
		Username: *username,
		Address:  *address,
	}
	status, err := hub.Status()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	found := false
	for key, light := range status.Lights {
		if light.Name != name {
			continue
		}
		found = true
		hub.ChangeLight(key).Transition(5).Brightness(int(255 * bri)).Send()
		break
	}
	if found {
		w.Write(ok)
	} else {
		w.Write([]byte("not found"))
	}
}

// call to signal that motion has been detected at the front door.
func motionAtEntry(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write(ok)
}

// call to signal all lights off.
// func allOff(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Content-Type", "text/plain")
// 	hub := &huego.Hub{
// 		Username: *username,
// 		Address:  *address,
// 	}
// 	status, err := hub.Status()
// 	if err != nil {
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	for key, light := range status.Lights {
// 		s := light.State
// 		s.On = false
// 		hub.SetLightState(key, s)
// 	}
// 	w.Write(ok)
// }

func main() {
	flag.Parse()
	http.HandleFunc("/lights/", setBrightness)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	err := http.ListenAndServe(":10443", nil)
	if err != nil {
		log.Fatal(err)
	}
}