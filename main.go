package main

/*
 * whatspeed is a tiny portable cli application written in go to test download/upload speeds
 */

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Config is
type Config struct {
	Client struct {
		Country string  `xml:"country,attr"`
		ISP     string  `xml:"isp,attr"`
		IP      string  `xml:"ip,attr"`
		Lat     float64 `xml:"lat,attr"`
		Lon     float64 `xml:"lon,attr"`
	} `xml:"client"`
}

// Servers is
type Servers struct {
	Servers []Server `xml:"servers>server"`
}

// Server is
type Server struct {
	URL     string  `xml:"url,attr"`
	Name    string  `xml:"name,attr"`
	Country string  `xml:"country,attr"`
	CC      string  `xml:"cc,attr"`
	Lat     float64 `xml:"lat,attr"`
	Lon     float64 `xml:"lon,attr"`
	ID      int     `xml:"id,attr"`
	Host    string  `xml:"host,attr"`
}

var (
	speedTestServerURLS = []string{
		// default to use tls but add option for http
		"https://www.speedtest.net/speedtest-servers-static.php",
		"https://c.speedtest.net/speedtest-servers-static.php",
		"https://www.speedtest.net/speedtest-servers.php",
		"https://c.speedtest.net/speedtest-servers.php",
	}
	speedTestConfigURL         = "https://www.speedtest.net/speedtest-config.php"
	earthRadius        float64 = 6378100
	downloadSizes              = []int{350, 500, 750, 1000, 1500, 2000, 2500}
)

func main() {
	var config *Config

	if err := xmlPls(speedTestConfigURL, &config); err != nil {
		log.Fatal(err)
	}

	servers := giveMeServers()
	nearestServer := nearestServerPls(config, servers)
	_, err := whatsMyDownloadSpeed(config, nearestServer)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// whatsMyDownloadSpeed will tell you your download speed, i hope its at least a gigabites
func whatsMyDownloadSpeed(config *Config, nearestSever *Server) (float64, error) {
	client := &http.Client{}
	var mbs float64

	for _, dlSize := range downloadSizes {
		splitURL := strings.Split(nearestSever.URL, "/")
		url := fmt.Sprintf("http:/%s/random%dx%d.jpg", strings.Join(splitURL[1:len(splitURL)-1], "/"), dlSize, dlSize)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return 0, err
		}
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("User-Agent", "whatSpeed")
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			return 0, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		mbs = (float64(len(body)*8) / float64(1000000)) / time.Now().Sub(start).Seconds()
		fmt.Println(mbs)
	}
	return mbs, nil
}

// xmlPls is a polite function that unmarshals your http response into xmlTarget
func xmlPls(url string, xmlTarget interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(bodyB, xmlTarget)
}

// nearestServerPls is a polite function that will find the nearest server, go figure?
// Uses Harversine formula to calculate the orthodromic distance between two server geo locations
func nearestServerPls(config *Config, servers []Server) *Server {
	smallestDistance := math.MaxFloat64
	var nServer Server
	for _, server := range servers {
		sLonR := server.Lon * math.Pi / 180
		sLatR := server.Lat * math.Pi / 180
		cLonR := config.Client.Lon * math.Pi / 180
		cLatR := config.Client.Lat * math.Pi / 180
		h := hsin(sLatR-cLatR) + math.Cos(cLatR)*math.Cos(sLatR)*hsin(sLonR-cLonR)
		dist := 2 * earthRadius * math.Asin(math.Sqrt(h))
		if dist <= smallestDistance {
			smallestDistance = dist
			nServer = server
		}
	}
	return &nServer
}

// giveMeServers all your server are belong to us
func giveMeServers() []Server {
	var serverWg sync.WaitGroup
	var servers []Server
	serversChan := make(chan *Servers, len(speedTestServerURLS))

	serverWg.Add(len(speedTestServerURLS))
	for _, serverURL := range speedTestServerURLS {
		go func(serverURL string, serverChan chan *Servers, wg *sync.WaitGroup) {
			defer wg.Done()
			var servers Servers
			if err := xmlPls(serverURL, &servers); err != nil {
				log.Fatal(err)
			}
			serverChan <- &servers
		}(serverURL, serversChan, &serverWg)
	}
	serverWg.Wait()
	close(serversChan)

	for serversFromChan := range serversChan {
		for _, server := range serversFromChan.Servers {
			servers = append(servers, server)
		}
	}
	return servers
}

// haversin function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
