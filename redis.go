package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

type RedisThief struct {
	Server        string
	Port          int
	ServerAddress string
	Connected     bool
	Timeout       int
	TargetType    string
	Connection    redis.Conn
	ServerInfo    json.RawMessage
}

// Connect tries to dial to the redis server
func (r *RedisThief) Connect(c chan Thief) {
	if r.ServerAddress == "" {
		r.ServerAddress = fmt.Sprintf("%s:%d", r.Server, r.Port)
	}

	rConn, err := redis.Dial("tcp", r.ServerAddress)

	if err != nil {
		logrus.Warn(err)
	} else {
		r.Connection = rConn
		r.Connected = true
	}

	c <- r
}

// GetServerInfo Pulls back all information known about the server
func (r *RedisThief) PullServerInfo(c chan Thief) {
	// If you're that kid -- we'll connect for ya
	if !r.Connected {
		panic("not connected")
	}
	result, err := redis.Bytes(r.Connection.Do("INFO"))

	if err != nil {
		logrus.Warn(err)
	}

	j := redisParseServerInfo(&result)
	r.ServerInfo = j

	c <- r

}

func (r RedisThief) GetServerInfo() json.RawMessage {
	return r.ServerInfo
}

// Close Closes the redis connection
func (r *RedisThief) Close() {
	r.Connection.Close()
	r.Connected = false
}

func (r RedisThief) GetTarget() string {
	return r.Server
}

func (r RedisThief) GetTargetType() string {
	return r.TargetType
}

func (r RedisThief) IsConnected() bool {
	return r.Connected
}

// This needs to be swapped to a lexer someday
func redisParseServerInfo(b *[]byte) json.RawMessage {
	// The map of values that we'll use to convert to JSON
	bannerMap := make(map[string]map[string]string)
	var parent string

	// Redis newlines on \r\n not just \n
	banner := strings.Split(string(*b), "\r\n")

	// For every line in the response
	// loop through and build our resulting JSON structure
	for _, line := range banner {

		// Get rid of empty lines
		if len(line) < 1 {
			continue
		}

		// Top levels of our JSON
		if string(line[0]) == "#" {
			// strip # and blank lines
			parent = string(line[2:])

			bannerMap[parent] = make(map[string]string)

			// Anything within that top level our JSON will now be a key value field
		} else {
			kv := strings.Split(line, ":")
			bannerMap[parent][kv[0]] = kv[1]
		}
	}

	jsonBanner, _ := json.Marshal(bannerMap)
	return jsonBanner
}
