package main

import (
	"./cinefade"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var action string
	flag.StringVar(&action, "action", "on", "lights on/off")
	flag.Parse()

	bridge := cinefade.GetBridge(false)

	switch action {
	case "on", "off":
		cinefade.ControlBulbs(bridge, action)
	case "register":
		cinefade.SaveBulbsState(bridge, "cinema.json")
	case "cinema":
		//TODO: must avoid action if bulbs are off
		_, err := os.Stat(cinefade.VarDir + "/cinema.json")
		if err != nil {
			log.Fatal("You must first use the register to save cinema lightstate")
		}
		cinefade.SaveBulbsState(bridge, "current.json")
		cinefade.SetBulbsState(bridge, "cinema.json")
	case "restore":
		_, err := os.Stat(cinefade.VarDir + "/current.json")
		if err != nil {
			log.Fatal("Can't call restore action without a backup state")
		}
		cinefade.SetBulbsState(bridge, "current.json")
	case "poll":
		client := cinefade.GetHttpClient()
		result := cinefade.Poll(client)
		fmt.Println(result)
	}
}
