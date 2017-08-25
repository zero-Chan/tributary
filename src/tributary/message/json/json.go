package json

import (
	"tributary/message"
)

type JsonMessage map[string]interface{}

func CreateJsonMessage() JsonMessage {
	v := make(JsonMessage)
	return v
}

func NewJsonMessage() *JsonMessage {
	v := CreateJsonMessage()

	return &v
}

func (msg JsonMessage) Type() string {
	return "JsonMessage"
}

func (msg JsonMessage) ADD(key string, val interface{}) {
	msg[key] = val
}

func (msg JsonMessage) GET(key string) interface{} {
	val, exist := msg[key]
	if !exist {
		return nil
	}

	return val
}

func (msg JsonMessage) Copy() message.Message {
	jm := CreateJsonMessage()
	for k, v := range msg {
		jm[k] = v
	}

	return jm
}
