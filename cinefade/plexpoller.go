package cinefade

import (
	"encoding/xml"
	"fmt"
	"github.com/mreiferson/go-httpclient"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	// URL = "http://muklo:32400/status/sessions"
	URL = "http://dockers:49158"
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

func poll(client *http.Client, c chan<- string) {
	for {
		req, _ := http.NewRequest("GET", URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("can't access plex", err)
		}
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
		time.Sleep(5000 * time.Millisecond)
	}
}

func CheckPlexStatus() {
	client := getHttpClient()
	c := make(chan string)
	go poll(client, c)
	for {
		result := <-c
		fmt.Println(result)
	}
}
