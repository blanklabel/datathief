package main

import (
	"encoding/json"
	"io/ioutil"

	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// TODO: Swap to new golang plugin system
// TODO: Fully concurrent -- connect == current, main app selects between getinfo, dump and connect
// TODO: Fully completed cmd interface
// TODO: Cassandra
// TODO: Elasticsearch
// TODO: Some datadumper

type Target struct {
	Server     string `json: server`
	Port       int    `json: port`
	TargetType string `json:"type"`
}

type Targets struct {
	Targets []Target `json:"targets"`
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("datathiefjson.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(file)
}

func main() {

	file, e := ioutil.ReadFile("./targets.json")
	if e != nil {
		logrus.Fatal("File error: %v\n", e)
	}

	t := Targets{}
	json.Unmarshal(file, &t)

	// Golang decides for ya that these are by reference
	connected := make(chan Thief)
	info := make(chan Thief)

	var wg sync.WaitGroup
	for _, target := range t.Targets {
		switch target.TargetType {

		case "redis":
			r := GetThief(REDIS, target.Server, target.Port)
			go r.Connect(connected)
			wg.Add(1)

		case "mongo":
			m := GetThief(MONGO, target.Server, target.Port)
			go m.Connect(connected)
			wg.Add(1)

		default:
			logrus.Warn("Unknown Target Type:", target.TargetType)
		}

	}

	// Wait for everything to finish and send the kill code when it does
	death := make(chan int)
	go func() {
		wg.Wait()
		death <- 0
	}()

	// Loop and wait for results
	for {
		select {
		case server := <-connected:

			if server.IsConnected() {
				go server.PullServerInfo(info)
			} else {
				wg.Done()
			}

		case results := <-info:
			defer results.Close()
			wg.Done()
			f := make(map[string]interface{})
			results.GetServerInfo()
			json.Unmarshal(results.GetServerInfo(), &f)
			switch results.GetTargetType() {

			case "REDIS":
				s := f["Server"].(map[string]interface{})
				logrus.WithFields(logrus.Fields{
					"target":  results.GetTarget(),
					"os":      s["os"],
					"version": s["redis_version"],
					"type":    results.GetTargetType(),
				}).Info()

			case "MONGO":
				logrus.WithFields(logrus.Fields{
					"target":   results.GetTarget(),
					"hostname": f["host"],
					"version":  f["version"],
					"type":     results.GetTargetType(),
				}).Info()

			}
		case code := <-death:
			os.Exit(code)
		}
	}
}
