package main

import (
	"flag"
	"fmt"

	"github.com/dichro/huego"
)

var (
	username = flag.String("username", "", "username for Hue hub")
	address = flag.String("address", "", "address of Hue hub")
)

func main() {
	flag.Parse()
	hub := &huego.Hub{
		Username: *username,
		Address: *address,
	}
	args := flag.Args()
	if len(args) == 1 && args[0] == "list" {
		list(hub, args)
		return
	}
	if len(args) == 3 && args[0] == "turn" && (args[1] != "on" || args[1] != "off") {
		turn(hub, args)
		return
	}
	usage(args)
}

func turn(hub *huego.Hub, args []string) {
	status, err := hub.Status()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	name := args[2]
	for key, light := range status.Lights {
		if light.Name == name {
			s := light.State
			s.On = args[1] == "on"
			hub.SetLightState(key, s)
		}
	}
}

func list(hub *huego.Hub, args []string) {
	status, err := hub.Status()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, light := range status.Lights {
		fmt.Println(light)
	}
}

func usage(args []string) {
	fmt.Println("bad", args)
}