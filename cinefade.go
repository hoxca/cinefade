package main

import (
	"./cinefade"
	"flag"
	"github.com/stretchr/goweb"
	"log"
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
	log.Printf("  visit: %s", Address)

	if listenErr != nil {
		log.Fatalf("Could not listen: %s", listenErr)
	}

	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			// stop the HTTP server
			log.Print("Stopping the server...")
			listener.Close()

			log.Print("Tearing down...")
			log.Fatal("Finished - bye bye.  ;-)")
		}
	}()

	log.Fatalf("Error in Serve: %s", s.Serve(listener))
}
