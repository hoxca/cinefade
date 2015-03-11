package cinefade

import (
	//	"errors"
	//	"fmt"
	"github.com/savaki/go.hue"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
	"log"
	//	"net"
	//"net/http"
	"os"
	//	"os/signal"
	// "strconv"
	// "time"
)

func cinefadeSwitch(bridge *hue.Bridge, action string) {
	switch action {
	case "on", "off":
		ControlBulbs(bridge, action)
	case "register":
		SaveBulbsState(bridge, "cinema.json")
	case "cinema":
		//TODO: must avoid action if bulbs are off
		_, err := os.Stat(VarDir + "/cinema.json")
		if err != nil {
			log.Fatal("You must first use the register to save cinema lightstate")
		}
		SaveBulbsState(bridge, "current.json")
		SetBulbsState(bridge, "cinema.json")
	case "restore":
		_, err := os.Stat(VarDir + "/current.json")
		if err != nil {
			log.Fatal("Can't call restore action without a backup state")
		}
		SetBulbsState(bridge, "current.json")
	case "poll":
		CheckPlexStatus()
	}
}

func MapRoutes(bridge *hue.Bridge) {
	goweb.MapBefore(func(c context.Context) error {
		// add a custom header
		c.HttpResponseWriter().Header().Set("X-Custom-Header", "Goweb")
		return nil
	})

	goweb.MapAfter(func(c context.Context) error {
		// TODO: log this
		log.Println("After resquest")
		return nil
	})

	goweb.Map("/", func(c context.Context) error {
		return goweb.Respond.With(c, 200, []byte("Welcome to the cinefade webapp"))
	})

	goweb.Map("/cinefade/{action}", func(c context.Context) error {
		// get the path value as an integer
		action := c.PathValue("action")
		cinefadeSwitch(bridge, action)
		// respond with the status
		return goweb.Respond.With(c, 200, []byte("action "+action+" was done by cinefade"))
	})

	goweb.Map(func(c context.Context) error {
		// just return a 404 message
		return goweb.API.Respond(c, 404, nil, []string{"File not found"})
	})
}
