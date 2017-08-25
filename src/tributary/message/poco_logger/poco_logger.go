package poco_logger

import (
	"fmt"
	"strconv"
	"time"

	"tributary/message"
)

type PocoLoggerMessage struct {
	loggerCustom
	Topic topic
	Tags  tags
	Log   string
}

func CreatePocoLoggerMessage() PocoLoggerMessage {
	v := PocoLoggerMessage{
		Tags: createTags(),
	}
	return v
}

func NewPocoLoggerMessage() *PocoLoggerMessage {
	v := CreatePocoLoggerMessage()
	return &v
}

func (msg PocoLoggerMessage) Type() string {
	return "PocoLoggerMessage"
}

func (msg *PocoLoggerMessage) ADD(key string, val interface{}) {
	msg.Tags.add(key, val)
}

func (msg *PocoLoggerMessage) GET(key string) interface{} {
	return msg.Tags.get(key)
}

func (msg PocoLoggerMessage) Copy() message.Message {
	m := msg
	return &m
}

type loggerCustom struct {
	Time    time.Time
	Level   string
	LineNum *int
}

type topic string

func (t topic) String() string {
	return string(t)
}

func (t *topic) SetString(s string) {
	*t = topic(s)
}

type tags map[string]string

func createTags() tags {
	t := make(tags)
	return t
}

func (t tags) add(key string, val interface{}) {
	switch val.(type) {
	case int:
		v := val.(int)
		t[key] = strconv.FormatInt(int64(v), 10)
	case int32:
		v := val.(int32)
		strconv.FormatInt(int64(v), 10)
		t[key] = strconv.FormatInt(int64(v), 10)
	case int64:
		v := val.(int64)
		t[key] = strconv.FormatInt(v, 10)
	case float32:
		v := val.(float32)
		t[key] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		v := val.(float64)
		t[key] = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		t[key] = strconv.FormatBool(val.(bool))
	case string:
		t[key] = val.(string)
	default:
		fmt.Printf("PocoLogger invalid ADD interface type[%T].", val)
		return
	}
}

func (t tags) get(key string) interface{} {
	v, exist := t[key]
	if !exist {
		return nil
	}

	return v
}
