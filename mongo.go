package main

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoThief class for interacting with mongodb
type MongoThief struct {
	Server        string
	Port          int
	ServerAddress string
	Connected     bool
	Timeout       int
	TargetType    string
	Connection    *mgo.Session
	ServerInfo    json.RawMessage
}

// Connect tries to dial to the mongo server
func (m *MongoThief) Connect(c chan Thief) {
	if m.ServerAddress == "" {
		m.ServerAddress = fmt.Sprintf("mongodb://%s:%d", m.Server, m.Port)
	}

	rConn, err := mgo.Dial(m.ServerAddress)

	if err != nil {
		logrus.Warn(err)
	} else {
		m.Connection = rConn
		m.Connected = true
	}

	c <- m
}

// PullServerInfo retrieves all information known about the server
func (m *MongoThief) PullServerInfo(c chan Thief) {
	// If you're that kid -- we'll connect for ya
	if !m.Connected {
		panic("not connected")
	}
	// This is really  map[string]interface{} so not relaly magic going on
	result := bson.M{}

	// Pulls all the information about the server
	err := m.Connection.Run(bson.D{{"serverStatus", 1}}, &result)
	if err != nil {
		logrus.Fatal(err)
	}
	// Add in database names
	result["databases"], err = m.Connection.DatabaseNames()

	// convert to JSON
	jsonBanner, err := json.Marshal(result)
	m.ServerInfo = jsonBanner

	c <- m
}

// GetServerInfo returns banner information
func (m MongoThief) GetServerInfo() json.RawMessage {
	return m.ServerInfo
}

// Close Closes the redis connection
func (m *MongoThief) Close() {
	m.Connection.Close()
	m.Connected = false
}

// GetTarget returns the target of the scan
func (m MongoThief) GetTarget() string {
	return m.Server
}

// GetTargetType returns the type of target (mongo, redis, etc)
func (m MongoThief) GetTargetType() string {
	return m.TargetType
}

// IsConnected returns a boolian if there is a connection active with the remote host
func (m MongoThief) IsConnected() bool {
	return m.Connected
}
