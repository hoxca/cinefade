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
	URL = "http://muklo:32400/status/sessions"
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

func GetHttpClient() *http.Client {
	transport := &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        4 * time.Second,
		ResponseHeaderTimeout: 2 * time.Second,
	}
	defer transport.Close()

	client := &http.Client{Transport: transport}
	return client
}

func Poll(client *http.Client) string {
	req, _ := http.NewRequest("GET", URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("can't access plex", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	//	body, err := ioutil.ReadFile("plex-stop.xml")
	var m MediaContainer
	xml.Unmarshal(body, &m)
	if err != nil {
		fmt.Printf("error: %v", err)
		return ""
	}
	switch m.Video.Player.State {
	case "":
		fmt.Println("Video is stopped")
		return "stopped"
	case "paused":
		fmt.Println("Video is paused")
		return "paused"
	case "playing":
		fmt.Println("Video is running")
		return "playing"
	default:
		return ""
	}
}
