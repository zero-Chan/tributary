package json

import (
	"encoding/json"
	//	"fmt"

	"tributary/message"
	msg_entity "tributary/message/json"
	"tributary/transcoder"
)

type JsonTransCoder struct {
}

func CreateJsonTransCoder() JsonTransCoder {
	v := JsonTransCoder{}

	return v
}

func NewJsonTransCoder() *JsonTransCoder {
	v := CreateJsonTransCoder()

	return &v
}

func (c *JsonTransCoder) New() transcoder.TransCoder {
	return NewJsonTransCoder()
}

func (c *JsonTransCoder) Marshal(msg message.Message) (data []byte, err error) {
	return json.Marshal(msg)
}

func (c *JsonTransCoder) Unmarshal(data []byte) (message.Message, error) {
	v := msg_entity.NewJsonMessage()

	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
