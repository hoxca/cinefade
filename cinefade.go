package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/savaki/go.hue"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	EtcDir = "/etc/cinefade"
	VarDir = "/var/lib/cinefade"
)

func GetAllBulbs(bridge *hue.Bridge) []*hue.Light {
	lights, err := bridge.GetAllLights()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return lights
}

func SaveBulbsState(bridge *hue.Bridge, filename string) {
	type bulbState struct {
		State hue.LightState `json:"state"`
		Name  string         `json:"name"`
	}
	var bulbs []*bulbState

	lights := GetAllBulbs(bridge)
	for _, light := range lights {
		bulbAttr, _ := light.GetLightAttributes()

		bulb := bulbState{
			Name:  bulbAttr.Name,
			State: bulbAttr.State,
		}

		bulbs = append(bulbs, &bulb)
	}

	jsonFile, err := os.Create(VarDir + "/" + filename)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer jsonFile.Close()

	bulbsState, _ := json.Marshal(bulbs)
	jsonFile.Write(bulbsState)
	jsonFile.Close()
}

func SetBulbsState(bridge *hue.Bridge, filename string) {
	type bulbState struct {
		State hue.LightState `json:"state"`
		Name  string         `json:"name"`
	}
	var bulbs []*bulbState

	lights := GetAllBulbs(bridge)

	bulbsState, err := ioutil.ReadFile(VarDir + "/" + filename)
	if err != nil {
		fmt.Println("Error:", err)
	}
	err = json.Unmarshal(bulbsState, &bulbs)
	if err != nil {
		fmt.Print("Error:", err)
	}
	for _, bulb := range bulbs {
		for _, light := range lights {
			if light.Name == bulb.Name {
				state := hue.SetLightState{
					On:     strconv.FormatBool(bulb.State.On),
					Bri:    strconv.Itoa(bulb.State.Bri),
					Hue:    strconv.Itoa(bulb.State.Hue),
					Sat:    strconv.Itoa(bulb.State.Sat),
					Xy:     bulb.State.Xy,
					Ct:     strconv.Itoa(bulb.State.Ct),
					Alert:  bulb.State.Alert,
					Effect: bulb.State.Effect,
				}
				light.SetState(state)
			}
		}
	}
}

func ControlBulbs(bridge *hue.Bridge, action string) {
	lights := GetAllBulbs(bridge)
	for _, light := range lights {
		if action == "off" {
			light.Off()
		} else {
			light.On()
			// light.SetLightTransition("80")
		}
	}
}

func main() {
	var action string
	flag.StringVar(&action, "action", "on", "lights on/off")
	flag.Parse()

	bridge := hue.NewBridge("192.168.1.99", "3d99dc627158727130a0d2a224445b")
	//bridge.Debug()

	switch action {
	case "on", "off":
		ControlBulbs(bridge, action)
	case "register":
		SaveBulbsState(bridge, "cinema.json")
	case "cinema":
		SaveBulbsState(bridge, "current.json")
		SetBulbsState(bridge, "cinema.json")
	case "restore":
		SetBulbsState(bridge, "current.json")
	}
}
