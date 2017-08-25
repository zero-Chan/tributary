package logger_alerter

import (
	"fmt"

	"tributary/message"
	"tributary/processor"
)

type LoggerAlerterProcessor struct {
	processor.NormalProcessor
}

func CreateLoggerAlerterProcessor() LoggerAlerterProcessor {
	v := LoggerAlerterProcessor{
		NormalProcessor: processor.CreateNormalProcessor(),
	}

	v.SetType("loggerAlerter")

	return v
}

func NewLoggerAlerterProcessor() *LoggerAlerterProcessor {
	v := CreateLoggerAlerterProcessor()

	return &v
}

func (p *LoggerAlerterProcessor) New() processor.Processor {
	return NewLoggerAlerterProcessor()
}

func (p *LoggerAlerterProcessor) StartUp() error {
	return nil
}

func (p *LoggerAlerterProcessor) Process(msg message.Message) (message.Message, error) {
	if msg == nil {
		return nil, fmt.Errorf("LoggerAlerterProcessor Get a nil msg.")
	}

	msg.ADD("log", "this is test txt_content.")

	return msg, nil
}
