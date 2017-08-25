package stdout

import (
	"fmt"
	"os"

	"tributary/message"
	"tributary/sink"
)

type StdoutSink struct {
	sink.NormalSink
	// rabbitmqProducer
}

func CreateStdoutSink() StdoutSink {
	v := StdoutSink{
		NormalSink: sink.CreateNormalSink(),
	}

	v.SetType("stdout")

	return v
}

func NewStdoutSink() *StdoutSink {
	v := CreateStdoutSink()
	return &v
}

func (s *StdoutSink) New() sink.Sink {
	return NewStdoutSink()
}

func (s *StdoutSink) StartUp() error {
	if s.Type() != "stdout" {
		return fmt.Errorf("StdoutSink StartUp fail: invalid type[%s].", s.Type())
	}

	err := s.NormalSink.StartUp()
	if err != nil {
		return err
	}

	return nil
}

func (s *StdoutSink) Publish(msg message.Message) error {
	data, err := s.Marshaler().Marshal(msg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stdout, "%s\n", string(data))
	if err != nil {
		return err
	}

	return nil
}
