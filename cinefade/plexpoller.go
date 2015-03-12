package cinefade

import (
	"encoding/xml"
	"github.com/mreiferson/go-httpclient"
	"github.com/savaki/go.hue"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	URL = "http://muklo:32400/status/sessions"
	//URL = "http://dockers:49160"
)

type Player struct {
	XMLName xml.Name `xml:"Player"`
	State   string   `xml:"state,attr"`
}

type Video struct {
	XMLName xml.Name `xml:"Video"`
	Player  Player   `xml:"Player"`
}

type MediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Size    int      `xml:"size,attr"`
	Video   Video    `xml:"Video"`
}

func getHttpClient() *http.Client {
	transport := &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        4 * time.Second,
		ResponseHeaderTimeout: 2 * time.Second,
	}
	defer transport.Close()

	client := &http.Client{Transport: transport}
	return client
}

func poll(client *http.Client, c chan string) {
	log.Println("Launch plex poller")
	for {
		req, _ := http.NewRequest("GET", URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("can't access plex", err)
		} else {
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			var m MediaContainer
			xml.Unmarshal(body, &m)
			if err != nil {
				log.Fatal("error: %v", err)
			}
			switch m.Video.Player.State {
			case "":
				c <- "stopped"
			case "paused":
				c <- "paused"
			case "playing":
				c <- "playing"
			default:
				c <- "unknown"
			}
		}
		time.Sleep(5000 * time.Millisecond)

	}
}

func hueControl(bridge *hue.Bridge, c chan string) {
	log.Println("Launch plex status")
	time.Sleep(1000 * time.Millisecond)
	var previous = "init"
	for {
		switch <-c {
		case "stopped", "paused":
			if previous == "cinema" || previous == "init" {
				log.Println("restore")
				cinefadeSwitch(bridge, "restore")
				previous = "stopped"
			}
		case "playing":
			if previous == "stopped" || previous == "init" {
				log.Println("cinema")
				cinefadeSwitch(bridge, "cinema")
				previous = "cinema"
			}
		default:
		}
	}
}

func RunCheckPlexStatus(bridge *hue.Bridge, c chan string) {
	client := getHttpClient()
	go poll(client, c)
	go hueControl(bridge, c)
}

func StopCheckPlexStatus() {
	log.Println("Close plex poller chan !")
	os.Exit(0)
}
