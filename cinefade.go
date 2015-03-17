package main

import (
	"./cinefade"
	"flag"
	"github.com/blackjack/syslog"
	"github.com/stretchr/goweb"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	Address string = ":9090"
)

func main() {

	syslog.Openlog("cinefade", syslog.LOG_PID, syslog.LOG_USER)

	var action string
	flag.StringVar(&action, "action", "on", "lights on/off")
	flag.Parse()

	bridge := cinefade.GetBridge(false)
	cinefade.MapRoutes(bridge)

	// make a http server using the goweb.DefaultHttpHandler()
	s := &http.Server{
		Addr:           Address,
		Handler:        goweb.DefaultHttpHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	listener, listenErr := net.Listen("tcp", Address)
	syslog.Info("Server cinefade started")
	syslog.Infof("visit: %s", Address)

	if listenErr != nil {
		syslog.Critf("Could not listen: %s", listenErr)
		os.Exit(1)
	}

	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			// stop the HTTP server
			syslog.Info("Stopping the cinefade server...")
			listener.Close()

			syslog.Info("Tearing down...")
			os.Exit(0)
		}
	}()

	syslog.Critf("Error in Serve: %s", s.Serve(listener))
	os.Exit(1)
}
