package main

import "encoding/json"

const (
	REDIS = iota
	MONGO
)

type Thief interface {
	Connect(chan Thief)
	PullServerInfo(chan Thief)
	GetServerInfo() json.RawMessage
	GetTarget() string
	IsConnected() bool
	GetTargetType() string
	Close()
}

func GetThief(thiefType int, server string, port int) Thief {
	switch thiefType {

	case REDIS:
		r := RedisThief{
			Server:     server,
			Port:       port,
			TargetType: "REDIS",
		}
		return &r

	case MONGO:
		m := MongoThief{
			Server:     server,
			Port:       port,
			TargetType: "MONGO",
		}
		return &m
	}

	// What did you ask me for? Because I'm pretty sure it doesn't exist
	return nil
}
