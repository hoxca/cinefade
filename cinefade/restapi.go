package cinefade

import (
	"github.com/blackjack/syslog"
	"github.com/savaki/go.hue"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
	"os"
	"time"
)

var r = make(chan bool)

func cinefadeExit() {
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
}

func cinefadeSwitch(bridge *hue.Bridge, action string) {

	c := make(chan string)
	switch action {
	case "on", "off":
		ControlBulbs(bridge, action)
	case "register":
		SaveBulbsState(bridge, "cinema.json")
	case "cinema":
		syslog.Info("Switch bulbs for cinema")
		if IsOneOfBulbsOn(bridge) {
			_, err := os.Stat(VarDir + "/cinema.json")
			if err != nil {
				syslog.Info("You must first use the register to save cinema lightstate")
			}
			SaveBulbsState(bridge, "current.json")
			SetBulbsState(bridge, "cinema.json")
		} else {
			syslog.Info("Bulbs are off, don't expect any change")
		}
	case "restore":
		syslog.Info("Restore bulbs settings")
		if IsOneOfBulbsOn(bridge) {
			_, err := os.Stat(VarDir + "/current.json")
			if err != nil {
				syslog.Info("Can't call restore action without a backup state")
			}
			SetBulbsState(bridge, "current.json")
		} else {
			syslog.Info("Bulbs are off, don't expect any change")
		}
	case "start":
		RunCheckPlexStatus(bridge, c)
	case "stop":
		StopCheckPlexStatus()
	case "exit":
		syslog.Info("Stopping the cinefade server...")
		go cinefadeExit()
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
