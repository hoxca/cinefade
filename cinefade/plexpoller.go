package cinefade

import (
	"encoding/xml"
	"github.com/ccding/go-config-reader/config"
	"github.com/mreiferson/go-httpclient"
	"github.com/savaki/go.hue"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	conf := config.NewConfig(EtcDir + "/cinefade.conf")
	err := conf.Read()
	if err != nil {
		log.Fatal(err)
	}
	plexUrl := conf.Get("", "plexUrl")

	for {
		select {
		case <-r:
			log.Println("! exit plex poller")
			return
		case <-time.After(5000 * time.Millisecond):
			log.Println("Exec plex poller")
			req, _ := http.NewRequest("GET", plexUrl, nil)
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
		}
	}
}

func hueControl(bridge *hue.Bridge, c chan string) {
	log.Println("Launch plex status")
	time.Sleep(1000 * time.Millisecond)
	var previous = "stopped"
	for {
		select {
		case <-r:
			log.Println("! exit plex status")
			return
		case <-time.After(5000 * time.Millisecond):
			log.Println("Exec status poller")
			switch <-c {
			case "stopped", "paused":
				if previous == "cinema" {
					log.Println("restore")
					cinefadeSwitch(bridge, "restore")
					previous = "stopped"
				}
			case "playing":
				if previous == "stopped" {
					log.Println("cinema")
					cinefadeSwitch(bridge, "cinema")
					previous = "cinema"
				}
			default:
			}
		}
	}
}

func RunCheckPlexStatus(bridge *hue.Bridge, c chan string) {
	client := getHttpClient()
	go poll(client, c)
	go hueControl(bridge, c)
}

func StopCheckPlexStatus() {
	r <- true
	r <- true
}
