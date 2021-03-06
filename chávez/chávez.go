package main

import (
	"errors"
	"flag"
	"fmt"
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
	re       = regexp.MustCompile("^/[a-z]+/([A-Za-z0-9' ]+)/([.0-9]+)(,([.0-9]+))?$")
	port     = flag.Int("port", 10443, "port to listen on")
)

func findLight(name string) (*huego.Change, error) {
	hub := &huego.Hub{
		Username: *username,
		Address:  *address,
	}
	status, err := hub.Status()
	if err != nil {
		return nil, err
	}
	if sw := status.FindSwitchable(name); sw != nil {
		return sw.Switch(), nil
	}
	return nil, errors.New("unknown light name")
}

func parseURL(req *http.Request) (c *huego.Change, vs []float64, err error) {
	groups := re.FindStringSubmatch(req.URL.Path)
	if len(groups) < 3 || len(groups) > 5 {
		log.Printf("Parse failed. Request %q parsed into %d groups.", req.URL.Path, len(groups))
		err = errors.New("parse error")
		return
	}
	name := groups[1]
	v, err := strconv.ParseFloat(groups[2], 64)
	if err != nil {
		log.Printf("Parse failed. Request %q has non-arg %q.", req.URL.Path, groups[2])
		return
	}
	vs = append(vs, v)
	if len(groups) > 4 && len(groups[4]) > 0 {
		v, err = strconv.ParseFloat(groups[4], 64)
		if err != nil {
			log.Printf("Parse failed. Request %q has non-arg %q.", req.URL.Path, groups[4])
			return
		}
		vs = append(vs, v)
	}
	c, err = findLight(name)
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

// wakeUp starts the morning wake-up mode.
func wakeUp(w http.ResponseWriter, req *http.Request) {
	change, err := findLight("Miki's Bedroom")
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		// turn on lights at mildest settings.
		change.Transition(0).State(true).Brightness(0).Temperature(500).Send()
		// now bring them up slowly.
		change.Reset().Transition(3000).Brightness(254).Temperature(154).Send()
		w.Write(ok)
	} else {
		w.Write([]byte(err.Error()))
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/brightness/", setBrightness)
	http.HandleFunc("/temperature/", setTemperature)
	http.HandleFunc("/state/", setState)
	http.HandleFunc("/colour/", setColour)
	http.HandleFunc("/wakeUp", wakeUp)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
