package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dichro/huego"
)

var (
	ok       = []byte("OK")
	username = flag.String("username", "", "username for Hue hub")
	address  = flag.String("address", "", "address of Hue hub")
	re       = regexp.MustCompile("^/[a-z]+/([A-Za-z0-9' ]+)/([.0-9]+)(,([.0-9]+))$")
)

func parseURL(req *http.Request) (c *huego.Change, vs []float64, err error) {
	groups := re.FindStringSubmatch(req.URL.Path)
	if len(groups) < 3 || len(groups) > 5 {
		err = errors.New("parse error")
		return
	}
	name := groups[1]
	v, err := strconv.ParseFloat(groups[2], 64)
	if err != nil {
		return
	}
	vs = append(vs, v)
	if len(groups) > 4 {
		v, err = strconv.ParseFloat(groups[4], 64)
		if err != nil {
			return
		}
		vs = append(vs, v)
	}
	hub := &huego.Hub{
		Username: *username,
		Address:  *address,
	}
	status, err := hub.Status()
	if err != nil {
		return
	}
	for key, light := range status.Lights {
		if light.Name == name {
			c = hub.ChangeLight(key)
			return
		}
	}
	err = errors.New("unknown light name")
	return
}

func setBrightness(w http.ResponseWriter, req *http.Request) {
	change, arg, err := parseURL(req)
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		change.Transition(5).Brightness(int(254 * arg[0])).Send()
		w.Write(ok)
	} else {
		w.Write([]byte(err.Error()))
	}
}

func setTemperature(w http.ResponseWriter, req *http.Request) {
	change, arg, err := parseURL(req)
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		change.Transition(5).Temperature(500 - int((500-154)*arg[0])).Send()
		w.Write(ok)
	} else {
		w.Write([]byte(err.Error()))
	}
}

func setState(w http.ResponseWriter, req *http.Request) {
	change, arg, err := parseURL(req)
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		change.State(arg[0] > 0.5).Send()
		w.Write(ok)
	} else {
		w.Write([]byte(err.Error()))
	}
}

func setColour(w http.ResponseWriter, req *http.Request) {
	change, arg, err := parseURL(req)
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		change.Transition(5).Colour(int(arg[0]*65535), int(arg[1]*254)).Send()
		w.Write(ok)
	} else {
		w.Write([]byte(err.Error()))
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
	http.HandleFunc("/brightness/", setBrightness)
	http.HandleFunc("/temperature/", setTemperature)
	http.HandleFunc("/state/", setState)
	http.HandleFunc("/colour/", setColour)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	err := http.ListenAndServe(":10443", nil)
	if err != nil {
		log.Fatal(err)
	}
}
