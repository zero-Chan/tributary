package sink

import (
	"tributary/message"
	"tributary/processor"
	"tributary/transcoder"
	"tributary/transcoder/factory"
)

type Productor interface {
	Publish([]byte) error
}

type Sink interface {
	New() Sink

	SetType(string)
	Type() string

	StartUp() error

	SetProcessor(processor.Processor)
	ProcessorKeys() []string

	SetMarshaler(transcoder.Marshaler)
	Marshaler() transcoder.Marshaler

	Publish(message.Message) error
}

type NormalSink struct {
	SinkType string `json:"type"`

	ProKeys          []string `json:"processors"`
	processorEititys []processor.Processor

	MarshalerKey    string `json:"marshaler"`
	marshalerEitity transcoder.Marshaler
}

func CreateNormalSink() NormalSink {
	v := NormalSink{
		ProKeys:          make([]string, 0),
		processorEititys: make([]processor.Processor, 0),
	}

	return v
}

func NewNormalSink() *NormalSink {
	v := CreateNormalSink()
	return &v
}

func (s *NormalSink) SetType(t string) {
	s.SinkType = t
}

func (s NormalSink) Type() string {
	return s.SinkType
}

func (s *NormalSink) SetProcessor(p processor.Processor) {
	s.processorEititys = append(s.processorEititys, p)
}

func (s NormalSink) ProcessorKeys() []string {
	return s.ProKeys
}

func (s *NormalSink) SetMarshaler(m transcoder.Marshaler) {
	s.marshalerEitity = m
}

func (s NormalSink) Marshaler() transcoder.Marshaler {
	return s.marshalerEitity
}

func (s *NormalSink) StartUp() error {
	tr, err := factory.CreateTranscoder(s.MarshalerKey)
	if err != nil {
		return err
	}

	s.SetMarshaler(tr)

	return nil
}
