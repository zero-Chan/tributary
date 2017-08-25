package poco_logger

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"tributary/message"
	msg_entity "tributary/message/poco_logger"
	"tributary/transcoder"
)

type PocoLoggerTransCoder struct {
}

func CreatePocoLoggerTransCoder() PocoLoggerTransCoder {
	v := PocoLoggerTransCoder{}
	return v
}

func NewPocoLoggerTransCoder() *PocoLoggerTransCoder {
	v := CreatePocoLoggerTransCoder()
	return &v
}

func (c *PocoLoggerTransCoder) New() transcoder.TransCoder {
	return NewPocoLoggerTransCoder()
}

func (c *PocoLoggerTransCoder) Marshal(msg message.Message) (data []byte, err error) {
	if msg == nil {
		return nil, fmt.Errorf("PocoLoggerTransCoder Marshal fail: msg is nil.")
	}

	pocoLog, ok := msg.(*msg_entity.PocoLoggerMessage)
	if !ok {
		return nil, fmt.Errorf("PocoLoggerTransCoder Marshal fail: invalid msg type[%T.]", msg)
	}

	// e.g : [2017-05-05 11:32:42.353] [INFO] [239] # [statistics/user/send_message] # [namespace=sns-circle] [task_id=0fa7be95-ff47-4d41-a49f-0197e66aec91] # {"sender":"10000","from":"client","peer":"150718211","to":"client","type":"announce","time":1493955162,"msg_id":2241809}
	logBuf := fmt.Sprintf("[%s] [%s] ", pocoLog.Time.Format("2006-01-02 15:04:05.999"), pocoLog.Level)
	if pocoLog.LineNum != nil {
		logBuf += fmt.Sprintf("[%d] ", *pocoLog.LineNum)
	}

	logBuf += fmt.Sprintf("# [%s] # ", pocoLog.Topic.String())

	var tagBuf string
	for key, val := range pocoLog.Tags {
		tagBuf += fmt.Sprintf("[%s=%s] ", key, val)
	}

	logBuf = logBuf + tagBuf + "# " + pocoLog.Log

	return []byte(logBuf), nil
}

func (c *PocoLoggerTransCoder) Unmarshal(data []byte) (message.Message, error) {
	if data == nil {
		return nil, fmt.Errorf("PocoLoggerTransCoder Unmarshal fail: data is nil.")
	}

	msg := msg_entity.NewPocoLoggerMessage()

	collect := strings.Split(string(data), " # ")
	if len(collect) != 4 {
		return nil, fmt.Errorf("Invalid PocoLogger format. %s", string(data))
	}

	logCustom := collect[0]
	topic := collect[1]
	tags := collect[2]
	logBuf := collect[3]

	var err error
	msg.Time, msg.Level, msg.LineNum, err = logCustomUnmarshal(logCustom)
	if err != nil {
		return nil, fmt.Errorf("PocoLoggerTransCoder Unmarshal fail: invalid custom params: %s", err)
	}

	topicContain := bracketsMatcher(topic)
	if len(topicContain) != 1 {
		return nil, fmt.Errorf("PocoLoggerTransCoder Unmarshal fail: invalid topic: %s", topic)
	}

	msg.Topic.SetString(topicContain[0])

	tagContainer, err := logTagsUnmarshal(tags)
	if err != nil {
		return nil, fmt.Errorf("PocoLoggerTransCoder Unmarshal fail: invalid tags params: %s", err)
	}

	for k, v := range tagContainer {
		msg.ADD(k, v)
	}

	msg.Log = logBuf

	return msg, nil
}

func bracketsMatcher(data string) []string {
	collect := make([]string, 0)

	posLeft := 0
	posRight := 0

	for {
		posLeft = strings.IndexByte(data[posRight:], '[')
		if posLeft == -1 {
			break
		}

		posLeft += posRight

		posRight = strings.IndexByte(data[posLeft:], ']')
		if posRight == -1 {
			break
		}
		posRight += posLeft

		collect = append(collect, data[posLeft+1:posRight])
	}

	return collect
}

func logCustomUnmarshal(data string) (t time.Time, level string, lineNum *int, err error) {
	collect := bracketsMatcher(data)

	if len(collect) != 2 && len(collect) != 3 {
		return time.Now(), "", nil, fmt.Errorf("Invalid format: %s", data)
	}

	t, err = time.Parse("2006-01-02 15:04:05.999", collect[0])
	if err != nil {
		return time.Now(), "", nil, fmt.Errorf("Parse time fail. err: %s .time: %s", err, collect[0])
	}

	level = collect[1]

	if len(collect) == 3 {
		lineNum = new(int)
		*lineNum, err = strconv.Atoi(collect[2])
		if err != nil {
			return time.Now(), "", nil, fmt.Errorf("Parse lineNum fail: %s", collect[2])
		}
	}

	return time.Now(), level, lineNum, nil
}

func logTagsUnmarshal(data string) (map[string]string, error) {
	container := make(map[string]string)

	collect := bracketsMatcher(data)
	for _, content := range collect {
		keyVal := strings.Split(content, "=")
		if len(keyVal) != 2 {
			continue
		}

		container[keyVal[0]] = keyVal[1]
	}

	return container, nil
}
