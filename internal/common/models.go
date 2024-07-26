package common

import "encoding/json"

type TPVUpdateMessage struct {
	MessageType           string      `json:"messageType"`
	MessageSequenceNumber json.Number `json:"messageSequenceNumber"`
	Device                Device      `json:"device"`
	TPV                   TPV         `json:"tpv"`
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TPV struct {
	Timestamp json.Number `json:"timestamp"`
	Position  Position    `json:"position"`
	Speed     json.Number `json:"speed"`
}

type Position struct {
	Latitude  json.Number `json:"latitude"`
	Longitude json.Number `json:"longitude"`
	Altitude  json.Number `json:"altitude"`
}

type FlatTPV struct {
	Timestamp int64
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
}
