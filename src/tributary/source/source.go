package source

import (
	"tributary/message"
	"tributary/processor"
	"tributary/transcoder"
	"tributary/transcoder/factory"
)

type Consumer interface {
	Accept() ([]byte, error)
}

type Source interface {
	New() Source

	SetType(string)
	Type() string

	StartUp() error

	SetProcessor(processor.Processor)
	ProcessorKeys() []string

	SetUnmarshaler(transcoder.Unmarshaler)
	Unmarshaler() transcoder.Unmarshaler

	Accept() (message.Message, error)
}

type NormalSource struct {
	SourceType string `json:"type"`

	PorKeys          []string `json:"processors"`
	processorEititys []processor.Processor

	UnmarshalerKey    string `json:"unmarshaler"`
	unmarshalerEitity transcoder.Unmarshaler

	Writer writer
}

func CreateNormalSource() NormalSource {
	v := NormalSource{
		PorKeys:          make([]string, 0),
		processorEititys: make([]processor.Processor, 0),
		Writer:           createWriter(0),
	}

	return v
}

func NewNormalSource() *NormalSource {
	v := CreateNormalSource()

	return &v
}

func (s *NormalSource) SetType(t string) {
	s.SourceType = t
}

func (s NormalSource) Type() string {
	return s.SourceType
}

func (s *NormalSource) SetProcessor(p processor.Processor) {
	s.processorEititys = append(s.processorEititys, p)
}

func (s NormalSource) ProcessorKeys() []string {
	return s.PorKeys
}

func (s *NormalSource) SetUnmarshaler(m transcoder.Unmarshaler) {
	s.unmarshalerEitity = m
}

func (s NormalSource) Unmarshaler() transcoder.Unmarshaler {
	return s.unmarshalerEitity
}

func (s *NormalSource) StartUp() error {
	tr, err := factory.CreateTranscoder(s.UnmarshalerKey)
	if err != nil {
		return err
	}

	s.SetUnmarshaler(tr)

	return nil
}

func (s *NormalSource) SetChannelSize(size int) {
	s.Writer.reset(size)
}

type writer struct {
	BufChannel chan []byte
}

func createWriter(size int) writer {
	w := writer{
		BufChannel: make(chan []byte, size),
	}

	return w
}

func newWriter(size int) *writer {
	w := createWriter(size)
	return &w
}

func (w *writer) reset(size int) {
	if w.BufChannel != nil {
		close(w.BufChannel)
	}

	w.BufChannel = make(chan []byte, size)
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.BufChannel <- p
	return len(p), nil
}
