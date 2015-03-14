package cinefade

import (
	"encoding/json"
	"fmt"
	"github.com/ccding/go-config-reader/config"
	"github.com/savaki/go.hue"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	EtcDir = "/etc/cinefade"
	VarDir = "/var/lib/cinefade"
)

func GetBridge(debug bool) *hue.Bridge {
	conf := config.NewConfig(EtcDir + "/cinefade.conf")
	err := conf.Read()
	if err != nil {
		log.Fatal(err)
	}

	hueIpAddr := conf.Get("", "hueIpAddr")
	hueUser := conf.Get("", "hueUser")

	if hueIpAddr == "" || hueUser == "" {
		log.Fatal("One of configuration hueIpAddr|hueUser not found")
	}

	bridge := hue.NewBridge(hueIpAddr, hueUser)
	if debug {
		bridge.Debug()
	}
	return bridge
}

func GetAllBulbs(bridge *hue.Bridge) []*hue.Light {
	lights, err := bridge.GetAllLights()
	if err != nil {
		fmt.Println("Error:", err)
	}
	return lights
}

func IsOneOfBulbsOn(bridge *hue.Bridge) bool {
	lights := GetAllBulbs(bridge)
	for _, light := range lights {
		bulbAttr, _ := light.GetLightAttributes()
		if bulbAttr.State.On {
			return true
		}
	}
	return false
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
		}
	}
}
