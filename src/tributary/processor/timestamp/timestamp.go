package timestamp

import (
	"fmt"
	"time"

	"tributary/message"
	"tributary/processor"
)

type TimestampProcessor struct {
	processor.NormalProcessor
}

func CreateTimestampProcessor() TimestampProcessor {
	v := TimestampProcessor{
		NormalProcessor: processor.CreateNormalProcessor(),
	}

	v.SetType("timestamp")

	return v
}

func NewTimestampProcessor() *TimestampProcessor {
	v := CreateTimestampProcessor()

	return &v
}

func (p *TimestampProcessor) New() processor.Processor {
	return NewTimestampProcessor()
}

func (p *TimestampProcessor) StartUp() error {
	return nil
}

func (p *TimestampProcessor) Process(msg message.Message) (message.Message, error) {
	if msg == nil {
		return nil, fmt.Errorf("TimestampProcessor Get a nil msg.")
	}

	msg.ADD("timestamp", time.Now().Unix())

	return msg, nil
}
